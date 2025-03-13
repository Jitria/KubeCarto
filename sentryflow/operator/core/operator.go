// SPDX-License-Identifier: Apache-2.0

package core

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"Operator/collector"
	"Operator/exporter"
)

// == //

// StopChan Channel
var StopChan chan struct{}

// init Function
func init() {
	StopChan = make(chan struct{})
}

// OperatorService Structure
type OperatorService struct {
	waitGroup *sync.WaitGroup
}

// NewOperator Function
func NewOperator() *OperatorService {
	sfo := new(OperatorService)
	sfo.waitGroup = new(sync.WaitGroup)
	return sfo
}

// DestroyOperator Function
func (sfo *OperatorService) DestroyOperator() {
	close(StopChan)

	// Stop collector
	if collector.StopCollector() {
		log.Print("[Operator] Stopped Collectors")
	} else {
		log.Print("[Operator] Failed to stop Collectors")
	}

	// Stop exporter
	if exporter.StopExporter() {
		log.Print("[Operator] Stopped Exporters")
	} else {
		log.Print("[Operator] Failed to stop Exporters")
	}

	log.Print("[Operator] Waiting for routine terminations")

	sfo.waitGroup.Wait()

	log.Print("[Operator] Terminated Operator")
}

// == //

// GetOSSigChannel Function
func GetOSSigChannel() chan os.Signal {
	c := make(chan os.Signal, 1)

	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		os.Interrupt)

	return c
}

// == //

// Operator Function
func Operator() {
	sfo := NewOperator()

	log.Print("[Operator] Initializing Operator")

	// == //

	// Start collector
	if !collector.StartCollector(sfo.waitGroup) {
		sfo.DestroyOperator()
		return
	}

	// Start exporter
	if !exporter.StartExporter(sfo.waitGroup) {
		sfo.DestroyOperator()
		return
	}

	log.Print("[Operator] Initialization is completed")

	// == //

	// listen for interrupt signals
	sigChan := GetOSSigChannel()
	<-sigChan
	log.Print("Got a signal to terminate Operator")

	// == //

	// Destroy Operator
	sfo.DestroyOperator()
}
