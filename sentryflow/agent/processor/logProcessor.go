// SPDX-License-Identifier: Apache-2.0

package processor

import (
	"log"
	"sync"

	"github.com/Jitria/SentryFlow/protobuf"

	"Agent/uploader"
)

// == //

// LogH global reference for Log Handler
var LogH *LogHandler

// init Function
func init() {
	LogH = NewLogHandler()
}

// LogHandler Structure
type LogHandler struct {
	stopChan chan struct{}

	apiLogChan  chan interface{}
	metricsChan chan interface{}
}

// NewLogHandler Structure
func NewLogHandler() *LogHandler {
	lh := &LogHandler{
		stopChan: make(chan struct{}),

		apiLogChan:  make(chan interface{}),
		metricsChan: make(chan interface{}),
	}

	return lh
}

// == //

// StartLogProcessor Function
func StartLogProcessor(wg *sync.WaitGroup) bool {
	// handle API logs
	go processAPILogs(wg)

	// handle Envoy metrics
	go processEnvoyMetrics(wg)

	log.Print("[LogProcessor] Started Log Processors")

	return true
}

// StopLogProcessor Function
func StopLogProcessor() bool {
	// One for processAPILogs
	LogH.stopChan <- struct{}{}

	// One for processMetrics
	LogH.stopChan <- struct{}{}

	log.Print("[LogProcessor] Stopped Log Processors")

	return true
}

// == //

// processAPILogs Function
func processAPILogs(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case logType, ok := <-LogH.apiLogChan:
			if !ok {
				log.Print("[LogProcessor] Failed to process an API log")
			}

			go uploader.UploadAPILog(logType.(*protobuf.APILog))

		case <-LogH.stopChan:
			wg.Done()
			return
		}
	}
}

// InsertAPILog Function
func InsertAPILog(data interface{}) {
	LogH.apiLogChan <- data
}

// processEnvoyMetrics Function
func processEnvoyMetrics(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case logType, ok := <-LogH.metricsChan:
			if !ok {
				log.Print("[LogProcessor] Failed to process Envoy metrics")
			}

			go uploader.UploadEnvoyMetrics(logType.(*protobuf.EnvoyMetrics))

		case <-LogH.stopChan:
			wg.Done()
			return
		}
	}
}

// InsertMetrics Function
func InsertMetrics(data interface{}) {
	LogH.metricsChan <- data
}

// == //
