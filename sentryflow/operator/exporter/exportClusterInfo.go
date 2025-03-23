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

// (2) ExpHandler에 slices와 채널 추가
// exportCluster.go 안에서만 쓰므로 ExpHandler 구조체를 확장할 수 있음.

func init() {
	// 보통 init()에서 ExpH 초기화하지만, 이미 다른 init()에 있다면 그쪽 사용
}

// InsertDeploy: 다른 모듈(collector/...)에서 "새로운 Deploy"가 왔을 때 호출
func InsertDeploy(dep *protobuf.Deploy) {
	// channel에 넣거나, 곧바로 ExpH.SendDeploy(dep) 호출 가능
	ExpH.exporterDeploy <- dep
}

// InsertPod: ...
func InsertPod(pod *protobuf.Pod) {
	ExpH.exporterPod <- pod
}

// InsertService: ...
func InsertService(svc *protobuf.Service) {
	ExpH.exporterSvc <- svc
}

// exportClusterHandler 함수: exportAPILogs 와 비슷하게 고루틴 돌며 채널 소비
func (exp *ExpHandler) exportClusterHandler(wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

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
			return
		}
	}
}

// (3) SendDeploy: 연결된 모든 deployStreamInform에게 dep 전송
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
		return errors.New(fmt.Sprintf("SendDeploy failed: %d/%d", failed, total))
	}
	return nil
}

// (4) SendPod, SendService 도 동일
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
		return errors.New(fmt.Sprintf("SendPod failed: %d/%d", failed, total))
	}
	return nil
}
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
		return errors.New(fmt.Sprintf("SendService failed: %d/%d", failed, total))
	}
	return nil
}

// (5) 실제 gRPC 핸들러: GetDeploy, GetPod, GetService
// ExpService 구조체에 메서드로 구현
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

	// 대기: 스트리밍은 서버가 일방적으로 계속 Send
	// 여기서는 client가 끊을 때까지 blocking
	return <-dsi.errChan
}

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
