// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"log"
	"mongo-client/mongodb"

	pb "github.com/Jitria/SentryFlow/protobuf"
)

// Feeder Structure
type Feeder struct {
	Running bool

	client            pb.SentryFlowClient
	logStream         pb.SentryFlow_GetAPILogClient
	envoyMetricStream pb.SentryFlow_GetEnvoyMetricsClient

	deployStream  pb.SentryFlow_GetDeployClient
	podStream     pb.SentryFlow_GetPodClient
	serviceStream pb.SentryFlow_GetServiceClient

	dbHandler mongodb.DBHandler

	Done chan struct{}
}

// NewClient Function
func NewClient(client pb.SentryFlowClient, clientInfo *pb.ClientInfo, logCfg string, metricCfg string, metricFilter string, mongoDBAddr string) *Feeder {
	fd := &Feeder{}

	fd.Running = true
	fd.client = client
	fd.Done = make(chan struct{})

	if logCfg != "none" {
		// Contact the server and print out its response
		logStream, err := client.GetAPILog(context.Background(), clientInfo)
		if err != nil {
			log.Fatalf("[Client] Could not get API log: %v", err)
		}

		fd.logStream = logStream
	}

	if metricCfg != "none" && (metricFilter == "all" || metricFilter == "envoy") {
		emStream, err := client.GetEnvoyMetrics(context.Background(), clientInfo)
		if err != nil {
			log.Fatalf("[Client] Could not get Envoy metrics: %v", err)
		}

		fd.envoyMetricStream = emStream
	}

	deployStr, err := client.GetDeploy(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get Deploy stream: %v", err)
	}
	fd.deployStream = deployStr

	podStr, err := client.GetPod(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get Pod stream: %v", err)
	}
	fd.podStream = podStr

	svcStr, err := client.GetService(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get Service stream: %v", err)
	}
	fd.serviceStream = svcStr

	// Initialize DB
	dbHandler, err := mongodb.NewMongoDBHandler(mongoDBAddr)
	if err != nil {
		log.Fatalf("[MongoDB] Unable to intialize DB: %v", err)
	}
	fd.dbHandler = *dbHandler

	return fd
}

// APILogRoutine Function
func (fd *Feeder) APILogRoutine(logCfg string) {
	for fd.Running {
		select {
		default:
			data, err := fd.logStream.Recv()
			if err != nil {
				log.Fatalf("[Client] Failed to receive an API log: %v", err)
				break
			}
			err = fd.dbHandler.InsertAPILog(data)
			if err != nil {
				log.Fatalf("[MongoDB] Failed to insert an API log: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// EnvoyMetricsRoutine Function
func (fd *Feeder) EnvoyMetricsRoutine(metricCfg string) {
	for fd.Running {
		select {
		default:
			data, err := fd.envoyMetricStream.Recv()
			if err != nil {
				log.Fatalf("[Client] Failed to receive Envoy metrics: %v", err)
				break
			}
			err = fd.dbHandler.InsertEnvoyMetrics(data)
			if err != nil {
				log.Fatalf("[MongoDB] Failed to insert Envoy metrics: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// DeployRoutine Function
func (fd *Feeder) DeployRoutine(clusterCfg string) {
	for fd.Running {
		select {
		case <-fd.Done:
			return
		default:
			dep, err := fd.deployStream.Recv()
			if err != nil {
				log.Printf("[Client] Deploy stream ended: %v", err)
				return
			}
			// MongoDB Insert
			if err := fd.dbHandler.InsertDeploy(dep); err != nil {
				log.Printf("[MongoDB] InsertDeploy error: %v", err)
			}
		}
	}
}

// PodRoutine Function
func (fd *Feeder) PodRoutine(clusterCfg string) {
	for fd.Running {
		select {
		case <-fd.Done:
			return
		default:
			pod, err := fd.podStream.Recv()
			if err != nil {
				log.Printf("[Client] Pod stream ended: %v", err)
				return
			}
			if err := fd.dbHandler.InsertPod(pod); err != nil {
				log.Printf("[MongoDB] InsertPod error: %v", err)
			}
		}
	}
}

// ServiceRoutine Function
func (fd *Feeder) ServiceRoutine(clusterCfg string) {
	for fd.Running {
		select {
		case <-fd.Done:
			return
		default:
			svc, err := fd.serviceStream.Recv()
			if err != nil {
				log.Printf("[Client] Service stream ended: %v", err)
				return
			}
			if err := fd.dbHandler.InsertService(svc); err != nil {
				log.Printf("[MongoDB] InsertService error: %v", err)
			}
		}
	}
}
