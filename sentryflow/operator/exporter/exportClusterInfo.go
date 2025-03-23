// SPDX-License-Identifier: Apache-2.0

package exporter

import (
	"fmt"
	"log"
	"sync"

	"github.com/Jitria/SentryFlow/protobuf"
)

// == //
type deployStreamInform struct {
	Hostname  string
	IPAddress string
	stream    protobuf.SentryFlow_GetDeployServer // stream Deploy
	errChan   chan error
}
type podStreamInform struct {
	Hostname  string
	IPAddress string
	stream    protobuf.SentryFlow_GetPodServer // stream Pod
	errChan   chan error
}
type svcStreamInform struct {
	Hostname  string
	IPAddress string
	stream    protobuf.SentryFlow_GetServiceServer // stream Service
	errChan   chan error
}

// InsertDeploy Function
func InsertDeploy(dep *protobuf.Deploy) {
	ExpH.exporterDeploy <- dep
}

// InsertPod Function
func InsertPod(pod *protobuf.Pod) {
	ExpH.exporterPod <- pod
}

// InsertService Function
func InsertService(svc *protobuf.Service) {
	ExpH.exporterSvc <- svc
}

// exportClusterHandler Function
func (exp *ExpHandler) exportClusterHandler(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case dep := <-exp.exporterDeploy:
			if dep != nil {
				exp.SendDeploy(dep)
			}
		case pod := <-exp.exporterPod:
			if pod != nil {
				exp.SendPod(pod)
			}
		case svc := <-exp.exporterSvc:
			if svc != nil {
				exp.SendService(svc)
			}

		case <-exp.stopChan:
			wg.Done()
			return
		}
	}
}

// SendDeploy Function
func (exp *ExpHandler) SendDeploy(dep *protobuf.Deploy) error {
	exp.exporterLock.Lock()
	defer exp.exporterLock.Unlock()

	failed := 0
	total := len(exp.deployExporters)
	for _, dsi := range exp.deployExporters {
		if err := dsi.stream.Send(dep); err != nil {
			failed++
			log.Printf("[Exporter] Failed to send Deploy to %s (%s): %v",
				dsi.Hostname, dsi.IPAddress, err)
		}
	}
	if failed > 0 {
		return fmt.Errorf("SendDeploy failed: %d/%d", failed, total)
	}
	return nil
}

// SendPod Function
func (exp *ExpHandler) SendPod(pod *protobuf.Pod) error {
	exp.exporterLock.Lock()
	defer exp.exporterLock.Unlock()

	failed := 0
	total := len(exp.podExporters)
	for _, psi := range exp.podExporters {
		if err := psi.stream.Send(pod); err != nil {
			failed++
			log.Printf("[Exporter] Failed to send Pod to %s (%s): %v",
				psi.Hostname, psi.IPAddress, err)
		}
	}
	if failed > 0 {
		return fmt.Errorf("SendPod failed: %d/%d", failed, total)
	}
	return nil
}

// SendService Function
func (exp *ExpHandler) SendService(svc *protobuf.Service) error {
	exp.exporterLock.Lock()
	defer exp.exporterLock.Unlock()

	failed := 0
	total := len(exp.svcExporters)
	for _, ssi := range exp.svcExporters {
		if err := ssi.stream.Send(svc); err != nil {
			failed++
			log.Printf("[Exporter] Failed to send Svc to %s (%s): %v",
				ssi.Hostname, ssi.IPAddress, err)
		}
	}
	if failed > 0 {
		return fmt.Errorf("SendService failed: %d/%d", failed, total)
	}
	return nil
}

// GetDeploy Function (for gRPC)
func (exs *ExpService) GetDeploy(info *protobuf.ClientInfo, stream protobuf.SentryFlow_GetDeployServer) error {
	log.Printf("[Exporter] Client %s (%s) connected to GetDeploy", info.HostName, info.IPAddress)

	dsi := &deployStreamInform{
		Hostname:  info.HostName,
		IPAddress: info.IPAddress,
		stream:    stream,
		errChan:   make(chan error),
	}

	ExpH.exporterLock.Lock()
	ExpH.deployExporters = append(ExpH.deployExporters, dsi)
	ExpH.exporterLock.Unlock()

	return <-dsi.errChan
}

// GetPod Function (for gRPC)
func (exs *ExpService) GetPod(info *protobuf.ClientInfo, stream protobuf.SentryFlow_GetPodServer) error {
	log.Printf("[Exporter] Client %s (%s) connected to GetPod", info.HostName, info.IPAddress)

	psi := &podStreamInform{
		Hostname:  info.HostName,
		IPAddress: info.IPAddress,
		stream:    stream,
		errChan:   make(chan error),
	}

	ExpH.exporterLock.Lock()
	ExpH.podExporters = append(ExpH.podExporters, psi)
	ExpH.exporterLock.Unlock()

	return <-psi.errChan
}

// GetService Function (for gRPC)
func (exs *ExpService) GetService(info *protobuf.ClientInfo, stream protobuf.SentryFlow_GetServiceServer) error {
	log.Printf("[Exporter] Client %s (%s) connected to GetService", info.HostName, info.IPAddress)

	ssi := &svcStreamInform{
		Hostname:  info.HostName,
		IPAddress: info.IPAddress,
		stream:    stream,
		errChan:   make(chan error),
	}

	ExpH.exporterLock.Lock()
	ExpH.svcExporters = append(ExpH.svcExporters, ssi)
	ExpH.exporterLock.Unlock()

	return <-ssi.errChan
}
