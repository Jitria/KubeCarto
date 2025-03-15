// SPDX-License-Identifier: Apache-2.0

package collector

import (
	"Operator/exporter"
	"Operator/protobuf"

	"fmt"
	"log"
	"net"
	"sync"

	"Operator/config"

	"google.golang.org/grpc"
)

type ColService struct {
	protobuf.UnimplementedOperatorServer
}

// == //

// ColH global reference for Collector Handler
var ColH *ColHandler

// ColHandler Structure
type ColHandler struct {
	colService  net.Listener
	grpcServer  *grpc.Server
	grpcService *ColService

	stopChan chan struct{}

	apiLogChan  chan interface{}
	metricsChan chan interface{}
}

// init Function
func init() {
	ColH = NewColHandler()
}

// NewColHandler Structure
func NewColHandler() *ColHandler {
	lh := &ColHandler{
		grpcService: new(ColService),

		stopChan: make(chan struct{}),

		apiLogChan:  make(chan interface{}),
		metricsChan: make(chan interface{}),
	}

	return lh
}

// == //

// StartCollector Function
func StartCollector(wg *sync.WaitGroup) bool {
	// Make a string with the given collector address and port
	collectorService := fmt.Sprintf("%s:%s", config.GlobalConfig.CollectorAddr, config.GlobalConfig.CollectorPort)

	// Start listening gRPC port
	colService, err := net.Listen("tcp", collectorService)
	if err != nil {
		log.Printf("[Collector] Failed to listen at %s: %v", collectorService, err)
		return false
	}
	ColH.colService = colService

	log.Printf("[Collector] Listening Collector gRPC services (%s)", collectorService)

	// Create gRPC Service
	gRPCServer := grpc.NewServer()
	ColH.grpcServer = gRPCServer

	protobuf.RegisterOperatorServer(gRPCServer, ColH.grpcService)

	// Serve gRPC Service
	go ColH.grpcServer.Serve(ColH.colService)

	log.Print("[Collector] Serving Collector gRPC services")

	// handle API logs
	go ProcessAPILogs(wg)

	// handle Envoy metrics
	go ProcessEnvoyMetrics(wg)

	log.Print("[LogProcessor] Started Log Processors")

	return true
}

// StopCollector Function
func StopCollector() bool {
	ColH.grpcServer.GracefulStop()

	log.Print("[Collector] Gracefully stopped Collector gRPC services")

	// One for ProcessAPILogs
	ColH.stopChan <- struct{}{}

	// One for ProcessMetrics
	ColH.stopChan <- struct{}{}

	log.Print("[LogProcessor] Stopped Log Processors")

	return true
}

// == //

// ProcessAPILogs Function
func ProcessAPILogs(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case logType, ok := <-ColH.apiLogChan:
			if !ok {
				log.Print("[LogProcessor] Failed to process an API log")
				continue
			}

			go exporter.InsertAPILog(logType.(*protobuf.APILog))

		case <-ColH.stopChan:
			wg.Done()
			return
		}
	}
}

// InsertAPILog Function
func InsertAPILog(data interface{}) {
	ColH.apiLogChan <- data
}

// ProcessEnvoyMetrics Function
func ProcessEnvoyMetrics(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case logType, ok := <-ColH.metricsChan:
			if !ok {
				log.Print("[LogProcessor] Failed to process Envoy metrics")
				continue
			}

			go exporter.InsertEnvoyMetrics(logType.(*protobuf.EnvoyMetrics))

		case <-ColH.stopChan:
			wg.Done()
			return
		}
	}
}

// InsertMetrics Function
func InsertMetrics(data interface{}) {
	ColH.metricsChan <- data
}
