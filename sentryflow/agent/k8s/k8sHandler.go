// SPDX-License-Identifier: Apache-2.0

package k8s

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"

	"k8s.io/apimachinery/pkg/util/json"

	"Agent/types"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

// == //

// K8sH global reference for Kubernetes Handler
var K8sH *KubernetesHandler

// init Function
func init() {
	K8sH = NewK8sHandler()
}

// KubernetesHandler Structure
type KubernetesHandler struct {
	config    *rest.Config
	clientSet *kubernetes.Clientset

	watchers  map[string]*cache.ListWatch
	informers map[string]cache.Controller

	podMap     map[string]*corev1.Pod        // NOT thread safe, key: Pod IP
	serviceMap map[string]*corev1.Service    // NOT thread safe, key: Service IP
	deployMap  map[string]*appsv1.Deployment // NOT thread safe, key: Namespace/DeploymentName
}

// NewK8sHandler Function
func NewK8sHandler() *KubernetesHandler {
	kh := &KubernetesHandler{
		watchers:  make(map[string]*cache.ListWatch),
		informers: make(map[string]cache.Controller),

		podMap:     make(map[string]*corev1.Pod),
		serviceMap: make(map[string]*corev1.Service),
		deployMap:  make(map[string]*appsv1.Deployment),
	}

	return kh
}

// == //

// InitK8sClient Function
func InitK8sClient() bool {
	var err error

	// Initialize in cluster config
	K8sH.config, err = rest.InClusterConfig()
	if err != nil {
		log.Print("[InitK8sClient] Failed to initialize Kubernetes client")
		return false
	}

	// Initialize Kubernetes clientSet
	K8sH.clientSet, err = kubernetes.NewForConfig(K8sH.config)
	if err != nil {
		log.Print("[InitK8sClient] Failed to initialize Kubernetes client")
		return false
	}

	// Create a mapping table for existing pods and services to IPs
	K8sH.initExistingResources()

	watchTargetsCoreV1 := []string{"pods", "services"}
	watchTargetsAppsV1 := []string{"deployments"}

	//  Initialize watchers for pods and services
	for _, target := range watchTargetsCoreV1 {
		watcher := cache.NewListWatchFromClient(
			K8sH.clientSet.CoreV1().RESTClient(),
			target,
			corev1.NamespaceAll,
			fields.Everything(),
		)
		K8sH.watchers[target] = watcher
	}

	// Initialize watchers for deployments
	for _, target := range watchTargetsAppsV1 {
		watcher := cache.NewListWatchFromClient(
			K8sH.clientSet.AppsV1().RESTClient(),
			target,
			corev1.NamespaceAll,
			fields.Everything(),
		)
		K8sH.watchers[target] = watcher
	}

	// Initialize informers
	K8sH.initInformers()

	log.Print("[InitK8sClient] Initialized Kubernetes client")

	return true
}

// initExistingResources Function that creates a mapping table for existing pods and services to IPs
// This is required since informers are NOT going to see existing resources until they are updated, created or deleted
// @todo: Refactor this function, this is kind of messy
func (k8s *KubernetesHandler) initExistingResources() {
	// List existing Pods
	podList, err := k8s.clientSet.CoreV1().Pods(corev1.NamespaceAll).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("[K8s] Failed to get Pods: %v", err.Error())
	} else {
		for _, pod := range podList.Items {
			currentPod := pod
			k8s.podMap[pod.Status.PodIP] = &currentPod
			log.Printf("[K8s] Add existing pod %s: %s/%s", pod.Status.PodIP, pod.Namespace, pod.Name)
		}
	}

	// List existing Services
	serviceList, err := k8s.clientSet.CoreV1().Services(corev1.NamespaceAll).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("[K8s] Failed to get Services: %v", err.Error())
	} else {
		for _, service := range serviceList.Items {
			currentService := service
			// Check if the service has a LoadBalancer type
			if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
				for _, lbIngress := range service.Status.LoadBalancer.Ingress {
					lbIP := lbIngress.IP
					if lbIP != "" {
						k8s.serviceMap[lbIP] = &currentService
						log.Printf("[K8s] Add existing service (LoadBalancer) %s: %s/%s", lbIP, service.Namespace, service.Name)
					}
				}
			} else {
				// ClusterIP
				if currentService.Spec.ClusterIP != "" && currentService.Spec.ClusterIP != "None" {
					k8s.serviceMap[currentService.Spec.ClusterIP] = &currentService
					log.Printf("[K8s] Add existing Service IP=%s -> %s/%s", currentService.Spec.ClusterIP, currentService.Namespace, currentService.Name)
				}
				// ExternalIPs
				for _, eip := range currentService.Spec.ExternalIPs {
					k8s.serviceMap[eip] = &currentService
					log.Printf("[K8s] Add existing Service ExternalIP=%s -> %s/%s", eip, currentService.Namespace, currentService.Name)
				}
			}
		}
	}

	// List existing Deployments
	deployList, err := k8s.clientSet.AppsV1().Deployments(corev1.NamespaceAll).List(context.TODO(), v1.ListOptions{})
	if err != nil {
		log.Printf("[K8s] Failed to get Deployments: %v", err)
	} else {
		for _, deploy := range deployList.Items {
			currentdeploy := deploy
			key := fmt.Sprintf("%s/%s", deploy.Namespace, deploy.Name)
			k8s.deployMap[key] = &currentdeploy
			log.Printf("[K8s] Add existing Deployment %s", key)
		}
	}
}

// initInformers Function that initializes informers for services and pods in a cluster
func (k8s *KubernetesHandler) initInformers() {
	// Create Pod controller informer
	if k8s.watchers["pods"] != nil {
		_, podController := cache.NewInformerWithOptions(
			cache.InformerOptions{
				ListerWatcher: k8s.watchers["pods"],
				ObjectType:    &corev1.Pod{},
				ResyncPeriod:  0,
				Handler: cache.ResourceEventHandlerFuncs{
					AddFunc: func(obj interface{}) {
						pod := obj.(*corev1.Pod)
						ip := pod.Status.PodIP
						if ip != "" {
							k8s.podMap[ip] = pod
						}
					},
					UpdateFunc: func(oldObj, newObj interface{}) {
						newPod := newObj.(*corev1.Pod)
						ip := newPod.Status.PodIP
						if ip != "" {
							k8s.podMap[ip] = newPod
						}
					},
					DeleteFunc: func(obj interface{}) {
						pod := obj.(*corev1.Pod)
						ip := pod.Status.PodIP
						if ip != "" {
							delete(k8s.podMap, ip)
						}
					},
				},
			},
		)
		k8s.informers["pods"] = podController
	}

	// Create Service controller informer
	if k8s.watchers["services"] != nil {
		_, svcController := cache.NewInformerWithOptions(
			cache.InformerOptions{
				ListerWatcher: k8s.watchers["services"],
				ObjectType:    &corev1.Service{},
				ResyncPeriod:  0, // 필요하면 조정
				Handler: cache.ResourceEventHandlerFuncs{
					AddFunc: func(obj interface{}) {
						svc := obj.(*corev1.Service)
						k8s.addOrUpdateServiceIPs(svc)
					},
					UpdateFunc: func(oldObj, newObj interface{}) {
						oldSvc := oldObj.(*corev1.Service)
						newSvc := newObj.(*corev1.Service)
						// 기존 IP 제거 후 새로 등록
						k8s.removeServiceIPs(oldSvc)
						k8s.addOrUpdateServiceIPs(newSvc)
					},
					DeleteFunc: func(obj interface{}) {
						svc := obj.(*corev1.Service)
						k8s.removeServiceIPs(svc)
					},
				},
			},
		)
		k8s.informers["services"] = svcController
	}

	// Create Deployment controller informer
	if k8s.watchers["deployments"] != nil {
		_, depController := cache.NewInformerWithOptions(
			cache.InformerOptions{
				ListerWatcher: k8s.watchers["deployments"],
				ObjectType:    &appsv1.Deployment{},
				ResyncPeriod:  0,
				Handler: cache.ResourceEventHandlerFuncs{
					AddFunc: func(obj interface{}) {
						dep := obj.(*appsv1.Deployment)
						key := fmt.Sprintf("%s/%s", dep.Namespace, dep.Name)
						k8s.deployMap[key] = dep
					},
					UpdateFunc: func(oldObj, newObj interface{}) {
						dep := newObj.(*appsv1.Deployment)
						key := fmt.Sprintf("%s/%s", dep.Namespace, dep.Name)
						k8s.deployMap[key] = dep
					},
					DeleteFunc: func(obj interface{}) {
						dep := obj.(*appsv1.Deployment)
						key := fmt.Sprintf("%s/%s", dep.Namespace, dep.Name)
						delete(k8s.deployMap, key)
					},
				},
			},
		)
		k8s.informers["deployments"] = depController
	}
}

// addOrUpdateServiceIPs Function
func (k8s *KubernetesHandler) addOrUpdateServiceIPs(svc *corev1.Service) {
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		for _, lbIngress := range svc.Status.LoadBalancer.Ingress {
			lbIP := lbIngress.IP
			if lbIP != "" {
				k8s.serviceMap[lbIP] = svc
			}
		}
	} else {
		if svc.Spec.ClusterIP != "" && svc.Spec.ClusterIP != "None" {
			k8s.serviceMap[svc.Spec.ClusterIP] = svc
		}
		for _, eip := range svc.Spec.ExternalIPs {
			k8s.serviceMap[eip] = svc
		}
	}
}

// removeServiceIPs Function
func (k8s *KubernetesHandler) removeServiceIPs(svc *corev1.Service) {
	if svc.Spec.Type == corev1.ServiceTypeLoadBalancer {
		for _, lbIngress := range svc.Status.LoadBalancer.Ingress {
			lbIP := lbIngress.IP
			if lbIP != "" {
				delete(k8s.serviceMap, lbIP)
			}
		}
	} else {
		if svc.Spec.ClusterIP != "" && svc.Spec.ClusterIP != "None" {
			delete(k8s.serviceMap, svc.Spec.ClusterIP)
		}
		for _, eip := range svc.Spec.ExternalIPs {
			delete(k8s.serviceMap, eip)
		}
	}
}

// == //

// RunInformers Function that starts running informers
func RunInformers(stopChan chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		for name, informer := range K8sH.informers {
			go func(name string, inf cache.Controller) {
				log.Printf("[RunInformers] Starting informer for %s", name)
				inf.Run(stopChan)
			}(name, informer)
		}
		<-stopChan
		log.Print("[RunInformers] Stop signal received. All informers stopping.")
	}()

	log.Print("[RunInformers] Started all Kubernetes informers")
}

// getConfigMap Function
func (k8s *KubernetesHandler) getConfigMap(namespace, name string) (string, error) {
	cm, err := k8s.clientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		log.Printf("[K8s] Failed to get ConfigMaps: %v", err)
		return "", err
	}

	// convert data to string
	data, err := json.Marshal(cm.Data)
	if err != nil {
		log.Printf("[K8s] Failed to marshal ConfigMap: %v", err)
		return "", err
	}

	return string(data), nil
}

// updateConfigMap Function
func (k8s *KubernetesHandler) updateConfigMap(namespace, name, data string) error {
	cm, err := k8s.clientSet.CoreV1().ConfigMaps(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		log.Printf("[K8s] Failed to get ConfigMap: %v", err)
		return err
	}

	if _, ok := cm.Data["mesh"]; !ok {
		return errors.New("[K8s] Unable to find field \"mesh\" from Istio config")
	}

	cm.Data["mesh"] = data
	if _, err := k8s.clientSet.CoreV1().ConfigMaps(namespace).Update(context.Background(), cm, v1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}

// PatchNamespaces Function that patches namespaces for adding 'istio-injection'
func PatchNamespaces() bool {
	namespaces, err := K8sH.clientSet.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Printf("[PatchNamespaces] Failed to get Namespaces: %v", err)
		return false
	}

	for _, ns := range namespaces.Items {
		namespace := ns.DeepCopy()

		// Skip the following namespaces
		// if namespace.Name == "sentryflow" {
		// 	continue
		// }

		namespace.Labels["istio-injection"] = "enabled"

		// Patch the namespace
		if _, err := K8sH.clientSet.CoreV1().Namespaces().Update(context.TODO(), namespace, v1.UpdateOptions{FieldManager: "patcher"}); err != nil {
			log.Printf("[PatchNamespaces] Failed to update Namespace %s: %v", namespace.Name, err)
			return false
		}

		log.Printf("[PatchNamespaces] Updated Namespace %s", namespace.Name)
	}

	log.Print("[PatchNamespaces] Updated all Namespaces")

	return true
}

// restartDeployment Function that performs a rolling restart for a deployment in the specified namespace
// @todo: fix this, this DOES NOT restart deployments
func (k8s *KubernetesHandler) restartDeployment(namespace string, deploymentName string) error {
	deploymentClient := k8s.clientSet.AppsV1().Deployments(namespace)

	// Get the deployment to retrieve the current spec
	deployment, err := deploymentClient.Get(context.Background(), deploymentName, v1.GetOptions{})
	if err != nil {
		return err
	}

	// Trigger a rolling restart by updating the deployment's labels or annotations
	deployment.Spec.Template.ObjectMeta.Labels["restartedAt"] = v1.Now().String()

	// Update the deployment to trigger the rolling restart
	_, err = deploymentClient.Update(context.TODO(), deployment, v1.UpdateOptions{})
	if err != nil {
		return err
	}

	return nil
}

// RestartDeployments Function that restarts the deployments in the namespaces with "istio-injection=enabled"
func RestartDeployments() bool {
	deployments, err := K8sH.clientSet.AppsV1().Deployments("").List(context.Background(), v1.ListOptions{})
	if err != nil {
		log.Printf("[PatchDeployments] Failed to get Deployments: %v", err)
		return false
	}

	for _, deployment := range deployments.Items {
		// Skip the following namespaces
		if deployment.Namespace == "sentryflow" {
			continue
		}

		// Restart the deployment
		if err := K8sH.restartDeployment(deployment.Namespace, deployment.Name); err != nil {
			log.Printf("[PatchDeployments] Failed to restart Deployment %s/%s: %v", deployment.Namespace, deployment.Name, err)
			return false
		}

		log.Printf("[PatchDeployments] Deployment %s/%s restarted", deployment.Namespace, deployment.Name)
	}

	log.Print("[PatchDeployments] Restarted all patched deployments")

	return true
}

// == //

// lookupIPAddress Function
func lookupIPAddress(ipAddr string) interface{} {
	// Look for pod map
	pod, ok := K8sH.podMap[ipAddr]
	if ok {
		return pod
	}

	// Look for service map
	service, ok := K8sH.serviceMap[ipAddr]
	if ok {
		return service
	}

	return nil
}

// LookupK8sResource Function
func LookupK8sResource(srcIP string) types.K8sResource {
	ret := types.K8sResource{
		Namespace: "Unknown",
		Name:      "Unknown",
		Labels:    make(map[string]string),
		Type:      types.K8sResourceTypeUnknown,
	}

	// Find Kubernetes resource from source IP (service or a pod)
	raw := lookupIPAddress(srcIP)

	// Currently supports Service or Pod
	switch raw.(type) {
	case *corev1.Pod:
		pod, ok := raw.(*corev1.Pod)
		if ok {
			ret.Namespace = pod.Namespace
			ret.Name = pod.Name
			ret.Labels = pod.Labels
			ret.Type = types.K8sResourceTypePod
		}
	case *corev1.Service:
		svc, ok := raw.(*corev1.Service)
		if ok {
			ret.Namespace = svc.Namespace
			ret.Name = svc.Name
			ret.Labels = svc.Labels
			ret.Type = types.K8sResourceTypeService
		}
	default:
		ret.Type = types.K8sResourceTypeUnknown
	}

	return ret
}

// == //
