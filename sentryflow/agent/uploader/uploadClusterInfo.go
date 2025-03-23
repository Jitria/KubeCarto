package uploader

import (
	"context"
	"log"
	"sync"
	"time"

	"Agent/config"
	"Agent/types"

	"github.com/Jitria/SentryFlow/protobuf"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// UploadClusterInfo Function
func UploadClusterInfo(clusterEvent *types.ClusterEvent) {
	UplH.clusterEvents <- clusterEvent
}

// uploadClusterInfo Function
func (upl *UplHandler) uploadClusterInfo(wg *sync.WaitGroup) {
	wg.Add(1)

	for {
		select {
		case evt := <-upl.clusterEvents:
			if evt == nil {
				continue
			}
			upl.handleClusterEvent(evt)

		case <-upl.stopChan:
			wg.Done()
			return
		}
	}
}

func (upl *UplHandler) handleClusterEvent(evt *types.ClusterEvent) {
	switch evt.ResourceType {
	case "Pod":
		upl.handlePodEvent(evt.Action, evt.Object)
	case "Service":
		upl.handleServiceEvent(evt.Action, evt.Object)
	case "Deploy":
		upl.handleDeployEvent(evt.Action, evt.Object)
	default:
		log.Printf("[Uploader] Unknown resource type: %s", evt.ResourceType)
	}
}

// == //

func (upl *UplHandler) handlePodEvent(action string, obj interface{}) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		log.Printf("[Uploader] handlePodEvent: Not a *corev1.Pod object")
		return
	}
	if upl.grpcClient == nil {
		return
	}

	// proto.Pod로 변환
	podProto := convertPodToProto(pod)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if action == "ADD" || action == "UPDATE" {
		// 없는 Update RPC 대신 AddPodInfo로 동일 처리
		resp, err := upl.grpcClient.AddPodInfo(ctx, podProto)
		if err != nil {
			log.Printf("[Uploader] Failed to AddPodInfo for Pod %s/%s: %v",
				pod.Namespace, pod.Name, err)
			return
		}
		log.Printf("[Uploader] handlePodEvent: %s Pod %s/%s => Operator resp=%v",
			action, pod.Namespace, pod.Name, resp)

	} else if action == "DELETE" {
		// Delete RPC
		delPodProto := &protobuf.Pod{
			Cluster:   podProto.Cluster,
			Namespace: pod.Namespace,
			Name:      pod.Name,
		}
		resp, err := upl.grpcClient.DeletePodInfo(ctx, delPodProto)
		if err != nil {
			log.Printf("[Uploader] Failed to DeletePodInfo for Pod %s/%s: %v",
				pod.Namespace, pod.Name, err)
			return
		}
		log.Printf("[Uploader] handlePodEvent: DELETE Pod %s/%s => Operator resp=%v",
			pod.Namespace, pod.Name, resp)
	}
}

func (upl *UplHandler) handleServiceEvent(action string, obj interface{}) {
	svc, ok := obj.(*corev1.Service)
	if !ok {
		log.Printf("[Uploader] handleServiceEvent: Not a *corev1.Service object")
		return
	}
	if upl.grpcClient == nil {
		return
	}

	svcProto := convertServiceToProto(svc)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if action == "ADD" || action == "UPDATE" {
		resp, err := upl.grpcClient.AddSvcInfo(ctx, svcProto)
		if err != nil {
			log.Printf("[Uploader] Failed to AddSvcInfo for Service %s/%s: %v",
				svc.Namespace, svc.Name, err)
			return
		}
		log.Printf("[Uploader] handleServiceEvent: %s Service %s/%s => Operator resp=%v",
			action, svc.Namespace, svc.Name, resp)

	} else if action == "DELETE" {
		delSvcProto := &protobuf.Service{
			Cluster:   svcProto.Cluster,
			Namespace: svc.Namespace,
			Name:      svc.Name,
		}
		resp, err := upl.grpcClient.DeleteSvcInfo(ctx, delSvcProto)
		if err != nil {
			log.Printf("[Uploader] Failed to DeleteSvcInfo for Service %s/%s: %v",
				svc.Namespace, svc.Name, err)
			return
		}
		log.Printf("[Uploader] handleServiceEvent: DELETE Service %s/%s => Operator resp=%v",
			svc.Namespace, svc.Name, resp)
	}
}

func (upl *UplHandler) handleDeployEvent(action string, obj interface{}) {
	dep, ok := obj.(*appsv1.Deployment)
	if !ok {
		log.Printf("[Uploader] handleDeployEvent: Not a *appsv1.Deployment object")
		return
	}
	if upl.grpcClient == nil {
		return
	}

	depProto := convertDeploymentToProto(dep)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if action == "ADD" || action == "UPDATE" {
		resp, err := upl.grpcClient.AddDeployInfo(ctx, depProto)
		if err != nil {
			log.Printf("[Uploader] Failed to AddDeployInfo for Deployment %s/%s: %v",
				dep.Namespace, dep.Name, err)
			return
		}
		log.Printf("[Uploader] handleDeployEvent: %s Deploy %s/%s => Operator resp=%v",
			action, dep.Namespace, dep.Name, resp)

	} else if action == "DELETE" {
		delDepProto := &protobuf.Deploy{
			Cluster:   depProto.Cluster,
			Namespace: dep.Namespace,
			Name:      dep.Name,
		}
		resp, err := upl.grpcClient.DeleteDeployInfo(ctx, delDepProto)
		if err != nil {
			log.Printf("[Uploader] Failed to DeleteDeployInfo for Deploy %s/%s: %v",
				dep.Namespace, dep.Name, err)
			return
		}
		log.Printf("[Uploader] handleDeployEvent: DELETE Deploy %s/%s => Operator resp=%v",
			dep.Namespace, dep.Name, resp)
	}
}

func convertPodToProto(pod *corev1.Pod) *protobuf.Pod {
	return &protobuf.Pod{
		Cluster:           config.GlobalConfig.ClusterName,
		Namespace:         pod.Namespace,
		Name:              pod.Name,
		NodeName:          pod.Spec.NodeName,
		PodIP:             pod.Status.PodIP,
		Status:            string(pod.Status.Phase),
		Labels:            pod.Labels,
		CreationTimestamp: pod.CreationTimestamp.String(),
	}
}

func convertServiceToProto(svc *corev1.Service) *protobuf.Service {
	var protoPorts []*protobuf.Port
	for _, p := range svc.Spec.Ports {
		protoPorts = append(protoPorts, &protobuf.Port{
			Port:       int32(p.Port),
			TargetPort: p.TargetPort.IntVal,
			Protocol:   string(p.Protocol),
		})
	}
	return &protobuf.Service{
		Cluster:   config.GlobalConfig.ClusterName,
		Namespace: svc.Namespace,
		Name:      svc.Name,
		Type:      string(svc.Spec.Type),
		ClusterIP: svc.Spec.ClusterIP,
		Ports:     protoPorts,
		Labels:    svc.Labels,
	}
}

func convertDeploymentToProto(dep *appsv1.Deployment) *protobuf.Deploy {
	replicas := int32(0)
	if dep.Spec.Replicas != nil {
		replicas = int32(*dep.Spec.Replicas)
	}
	return &protobuf.Deploy{
		Cluster:           config.GlobalConfig.ClusterName,
		Namespace:         dep.Namespace,
		Name:              dep.Name,
		DesiredReplicas:   replicas,
		AvailableReplicas: dep.Status.AvailableReplicas,
		Labels:            dep.Labels,
		CreationTimestamp: dep.CreationTimestamp.String(),
	}
}
