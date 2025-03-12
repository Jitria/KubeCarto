// SPDX-License-Identifier: Apache-2.0

package core

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"sentryflow/collector"
	"sentryflow/config"
	"sentryflow/exporter"
	"sentryflow/k8s"
	"sentryflow/processor"
)

// == //

// StopChan Channel
var StopChan chan struct{}

// init Function
func init() {
	StopChan = make(chan struct{})
}

// SentryFlowService Structure
type SentryFlowService struct {
	waitGroup *sync.WaitGroup
}

// NewSentryFlow Function
func NewSentryFlow() *SentryFlowService {
	sf := new(SentryFlowService)
	sf.waitGroup = new(sync.WaitGroup)
	return sf
}

// DestroySentryFlow Function
func (sf *SentryFlowService) DestroySentryFlow() {
	close(StopChan)

	// Remove SentryFlow collector config from Kubernetes
	if k8s.UnpatchIstioConfigMap() {
		log.Print("[SentryFlow] Unpatched Istio ConfigMap")
	} else {
		log.Print("[SentryFlow] Failed to unpatch Istio ConfigMap")
	}

	// Stop collector
	if collector.StopCollector() {
		log.Print("[SentryFlow] Stopped Collectors")
	} else {
		log.Print("[SentryFlow] Failed to stop Collectors")
	}

	// Stop Log Processor
	if processor.StopLogProcessor() {
		log.Print("[SentryFlow] Stopped Log Processors")
	} else {
		log.Print("[SentryFlow] Failed to stop Log Processors")
	}

	// Stop exporter
	if exporter.StopExporter() {
		log.Print("[SentryFlow] Stopped Exporters")
	} else {
		log.Print("[SentryFlow] Failed to stop Exporters")
	}

	log.Print("[SentryFlow] Waiting for routine terminations")

	sf.waitGroup.Wait()

	log.Print("[SentryFlow] Terminated SentryFlow")
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

// SentryFlow Function
func SentryFlow() {
	sf := NewSentryFlow()

	log.Print("[SentryFlow] Initializing SentryFlow")

	// == //

	// Initialize Kubernetes client
	if !k8s.InitK8sClient() {
		sf.DestroySentryFlow()
		return
	}

	// Start Kubernetes informers
	k8s.RunInformers(StopChan, sf.waitGroup)

	// Patch Istio ConfigMap
	if !k8s.PatchIstioConfigMap() {
		sf.DestroySentryFlow()
		return
	}

	// Patch Namespaces
	if config.GlobalConfig.PatchingNamespaces {
		if !k8s.PatchNamespaces() {
			sf.DestroySentryFlow()
			return
		}
	}

	// Patch Deployments
	if config.GlobalConfig.RestartingPatchedDeployments {
		if !k8s.RestartDeployments() {
			sf.DestroySentryFlow()
			return
		}
	}

	// == //

	// Start collector
	if !collector.StartCollector() {
		sf.DestroySentryFlow()
		return
	}

	// Start log processor
	if !processor.StartLogProcessor(sf.waitGroup) {
		sf.DestroySentryFlow()
		return
	}

	// Start exporter
	if !exporter.StartExporter(sf.waitGroup) {
		sf.DestroySentryFlow()
		return
	}

	log.Print("[SentryFlow] Initialization is completed")

	// == //

	// listen for interrupt signals
	sigChan := GetOSSigChannel()
	<-sigChan
	log.Print("Got a signal to terminate SentryFlow")

	// == //

	// Destroy SentryFlow
	sf.DestroySentryFlow()
}
