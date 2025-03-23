// SPDX-License-Identifier: Apache-2.0

package exporter

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"github.com/Jitria/SentryFlow/protobuf"
)

// == //

// envoyMetricsStreamInform structure
type envoyMetricsStreamInform struct {
	Hostname  string
	IPAddress string

	metricsStream protobuf.SentryFlow_GetEnvoyMetricsServer

	error chan error
}

// InsertEnvoyMetrics Function
func InsertEnvoyMetrics(evyMetrics *protobuf.EnvoyMetrics) {
	ExpH.exporterMetrics <- evyMetrics
}

// exportEnvoyMetrics Function
func (exp *ExpHandler) exportEnvoyMetrics(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case evyMetrics, ok := <-exp.exporterMetrics:
			if !ok {
				log.Printf("[Exporter] Failed to fetch metrics from Envoy Metrics channel")
				wg.Done()
				return
			}

			if err := exp.SendEnvoyMetrics(evyMetrics); err != nil {
				log.Printf("[Exporter] Failed to export Envoy metrics: %v", err)
			}

		case <-exp.stopChan:
			wg.Done()
			return
		}
	}
}

// SendEnvoyMetrics Function
func (exp *ExpHandler) SendEnvoyMetrics(evyMetrics *protobuf.EnvoyMetrics) error {
	failed := 0
	total := len(exp.envoyMetricsExporters)

	for _, exporter := range exp.envoyMetricsExporters {
		if err := exporter.metricsStream.Send(evyMetrics); err != nil {
			log.Printf("[Exporter] Failed to export Envoy metrics to %s(%s): %v", exporter.Hostname, exporter.IPAddress, err)
			failed++
		}
	}

	if failed != 0 {
		msg := fmt.Sprintf("[Exporter] Failed to export Envoy metrics properly (%d/%d failed)", failed, total)
		return errors.New(msg)
	}

	return nil
}

// GetEnvoyMetrics Function (for gRPC)
func (exs *ExpService) GetEnvoyMetrics(info *protobuf.ClientInfo, stream protobuf.SentryFlow_GetEnvoyMetricsServer) error {
	log.Printf("[Exporter] Client %s (%s) connected (GetEnvoyMetrics)", info.HostName, info.IPAddress)

	currExporter := &envoyMetricsStreamInform{
		Hostname:      info.HostName,
		IPAddress:     info.IPAddress,
		metricsStream: stream,
	}

	ExpH.exporterLock.Lock()
	ExpH.envoyMetricsExporters = append(ExpH.envoyMetricsExporters, currExporter)
	ExpH.exporterLock.Unlock()

	return <-currExporter.error
}

// == //
