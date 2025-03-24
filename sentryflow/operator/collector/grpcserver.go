// SPDX-License-Identifier: Apache-2.0

package collector

import (
	"Operator/exporter"
	"context"
	"fmt"
	"io"
	"log"
	"sync"

	"github.com/Jitria/SentryFlow/protobuf"
)

//////////////////
// ClusterEvent //
//////////////////

// Deploy Function
func (cs *ColService) AddDeployEvent(ctx context.Context, dep *protobuf.Deploy) (*protobuf.Response, error) {
	log.Printf("[Operator] AddDeployEvent: got Deploy %s/%s cluster=%s", dep.Namespace, dep.Name, dep.Cluster)

	exporter.InsertDeployAdd(dep)
	return &protobuf.Response{Msg: 0}, nil
}

func (cs *ColService) UpdateDeployEvent(ctx context.Context, dep *protobuf.Deploy) (*protobuf.Response, error) {
	log.Printf("[Operator] UpdateDeployEvent: got Deploy %s/%s cluster=%s", dep.Namespace, dep.Name, dep.Cluster)

	exporter.InsertDeployUpdate(dep)
	return &protobuf.Response{Msg: 0}, nil
}

func (cs *ColService) DeleteDeployEvent(ctx context.Context, dep *protobuf.Deploy) (*protobuf.Response, error) {
	log.Printf("[Operator] DeleteDeployEvent: got Deploy %s/%s cluster=%s", dep.Namespace, dep.Name, dep.Cluster)

	exporter.InsertDeployDelete(dep)
	return &protobuf.Response{Msg: 0}, nil
}

// Pod Function
func (cs *ColService) AddPodEvent(ctx context.Context, pod *protobuf.Pod) (*protobuf.Response, error) {
	log.Printf("[Operator] AddPodEvent: got Pod %s/%s cluster=%s IP=%s",
		pod.Namespace, pod.Name, pod.Cluster, pod.PodIP)

	exporter.InsertPodAdd(pod)
	return &protobuf.Response{Msg: 0}, nil
}

func (cs *ColService) UpdatePodEvent(ctx context.Context, pod *protobuf.Pod) (*protobuf.Response, error) {
	log.Printf("[Operator] UpdatePodEvent: got Pod %s/%s cluster=%s IP=%s",
		pod.Namespace, pod.Name, pod.Cluster, pod.PodIP)

	exporter.InsertPodUpdate(pod)
	return &protobuf.Response{Msg: 0}, nil
}

func (cs *ColService) DeletePodEvent(ctx context.Context, pod *protobuf.Pod) (*protobuf.Response, error) {
	log.Printf("[Operator] DeletePodEvent: got Pod %s/%s cluster=%s",
		pod.Namespace, pod.Name, pod.Cluster)

	exporter.InsertPodDelete(pod)
	return &protobuf.Response{Msg: 0}, nil
}

// Service Function
func (cs *ColService) AddSvcEvent(ctx context.Context, svc *protobuf.Service) (*protobuf.Response, error) {
	log.Printf("[Operator] AddSvcEvent: got Service %s/%s cluster=%s clusterIP=%s",
		svc.Namespace, svc.Name, svc.Cluster, svc.ClusterIP)

	exporter.InsertSvcAdd(svc)
	return &protobuf.Response{Msg: 0}, nil
}

func (cs *ColService) UpdateSvcEvent(ctx context.Context, svc *protobuf.Service) (*protobuf.Response, error) {
	log.Printf("[Operator] UpdateSvcEvent: got Service %s/%s cluster=%s clusterIP=%s",
		svc.Namespace, svc.Name, svc.Cluster, svc.ClusterIP)

	exporter.InsertSvcUpdate(svc)
	return &protobuf.Response{Msg: 0}, nil
}

func (cs *ColService) DeleteSvcEvent(ctx context.Context, svc *protobuf.Service) (*protobuf.Response, error) {
	log.Printf("[Operator] DeleteSvcEvent: got Service %s/%s cluster=%s",
		svc.Namespace, svc.Name, svc.Cluster)

	exporter.InsertSvcDelete(svc)
	return &protobuf.Response{Msg: 0}, nil
}

////////////
// APILog //
////////////

// GiveAPILog Function
func (cs *ColService) GiveAPILog(stream protobuf.SentryFlow_GiveAPILogServer) error {
	for {
		// Receive APILog from stream.
		apiLog, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&protobuf.Response{Msg: 0})
		}
		if err != nil {
			return fmt.Errorf("GiveAPILog recv error: %v", err)
		}
		ColH.apiLogChan <- apiLog
	}
}

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

/////////////////
// EnovyMetric //
/////////////////

// GiveEnvoyMetrics Function
func (cs *ColService) GiveEnvoyMetrics(stream protobuf.SentryFlow_GiveEnvoyMetricsServer) error {
	for {
		// Receive EnvoyMetrics from stream.
		envoyMetrics, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&protobuf.Response{Msg: 0})
		}
		if err != nil {
			return fmt.Errorf("GiveEnvoyMetrics recv error: %v", err)
		}
		ColH.metricsChan <- envoyMetrics
	}
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

// == //
