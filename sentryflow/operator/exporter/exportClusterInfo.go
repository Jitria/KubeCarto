// SPDX-License-Identifier: Apache-2.0

package exporter

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"

	"github.com/Jitria/SentryFlow/protobuf"
)

// == //

// apiLogStreamInform structure
type apiLogStreamInform struct {
	Hostname  string
	IPAddress string

	stream protobuf.SentryFlow_GetAPILogServer

	error chan error
}

// InsertAPILog Function
func InsertAPILog(apiLog *protobuf.APILog) {
	ExpH.exporterAPILogs <- apiLog

	// Make a string with labels
	var labelString []string
	for k, v := range apiLog.SrcLabel {
		labelString = append(labelString, fmt.Sprintf("%s:%s", k, v))
	}
	sort.Strings(labelString)
}

// exportAPILogs Function
func (exp *ExpHandler) exportAPILogs(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case apiLog, ok := <-exp.exporterAPILogs:
			if !ok {
				log.Printf("[Exporter] Failed to fetch APIs from APIs channel")
				wg.Done()
				return
			}

			if err := exp.SendAPILogs(apiLog); err != nil {
				log.Printf("[Exporter] Failed to export API Logs: %v", err)
			}

		case <-exp.stopChan:
			wg.Done()
			return
		}
	}
}

// SendAPILogs Function
func (exp *ExpHandler) SendAPILogs(apiLog *protobuf.APILog) error {
	failed := 0
	total := len(exp.apiLogExporters)

	for _, exporter := range exp.apiLogExporters {
		if err := exporter.stream.Send(apiLog); err != nil {
			log.Printf("[Exporter] Failed to export an API log to %s (%s): %v", exporter.Hostname, exporter.IPAddress, err)
			failed++
		}
	}

	if failed != 0 {
		msg := fmt.Sprintf("[Exporter] Failed to export API logs properly (%d/%d failed)", failed, total)
		return errors.New(msg)
	}

	return nil
}

// GetAPILog Function (for gRPC)
func (exs *ExpService) GetAPILog(info *protobuf.ClientInfo, stream protobuf.SentryFlow_GetAPILogServer) error {
	log.Printf("[Exporter] Client %s (%s) connected (GetAPILog)", info.HostName, info.IPAddress)

	currExporter := &apiLogStreamInform{
		Hostname:  info.HostName,
		IPAddress: info.IPAddress,
		stream:    stream,
	}

	ExpH.exporterLock.Lock()
	ExpH.apiLogExporters = append(ExpH.apiLogExporters, currExporter)
	ExpH.exporterLock.Unlock()

	return <-currExporter.error
}

// == //
