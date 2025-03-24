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

	deployAddStream    pb.SentryFlow_AddDeployEventDBClient
	deployUpdateStream pb.SentryFlow_UpdateDeployEventDBClient
	deployDeleteStream pb.SentryFlow_DeleteDeployEventDBClient

	podAddStream    pb.SentryFlow_AddPodEventDBClient
	podUpdateStream pb.SentryFlow_UpdatePodEventDBClient
	podDeleteStream pb.SentryFlow_DeletePodEventDBClient

	svcAddStream    pb.SentryFlow_AddSvcEventDBClient
	svcUpdateStream pb.SentryFlow_UpdateSvcEventDBClient
	svcDeleteStream pb.SentryFlow_DeleteSvcEventDBClient

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

	// Deploy Add
	addDepStr, err := client.AddDeployEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get AddDeployEventDB stream: %v", err)
	}
	fd.deployAddStream = addDepStr

	// Deploy Update
	updDepStr, err := client.UpdateDeployEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get UpdateDeployEventDB stream: %v", err)
	}
	fd.deployUpdateStream = updDepStr

	// Deploy Delete
	delDepStr, err := client.DeleteDeployEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get DeleteDeployEventDB stream: %v", err)
	}
	fd.deployDeleteStream = delDepStr

	// Pod Add
	addPodStr, err := client.AddPodEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get AddPodEventDB stream: %v", err)
	}
	fd.podAddStream = addPodStr

	// Pod Update
	updPodStr, err := client.UpdatePodEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get UpdatePodEventDB stream: %v", err)
	}
	fd.podUpdateStream = updPodStr

	// Pod Delete
	delPodStr, err := client.DeletePodEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get DeletePodEventDB stream: %v", err)
	}
	fd.podDeleteStream = delPodStr

	// Service Add
	addSvcStr, err := client.AddSvcEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get AddSvcEventDB stream: %v", err)
	}
	fd.svcAddStream = addSvcStr

	// Service Update
	updSvcStr, err := client.UpdateSvcEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get UpdateSvcEventDB stream: %v", err)
	}
	fd.svcUpdateStream = updSvcStr

	// Service Delete
	delSvcStr, err := client.DeleteSvcEventDB(context.Background(), clientInfo)
	if err != nil {
		log.Fatalf("[Client] Could not get DeleteSvcEventDB stream: %v", err)
	}
	fd.svcDeleteStream = delSvcStr

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

// DeployAddRoutine Function
func (fd *Feeder) DeployAddRoutine() {
	for fd.Running {
		select {
		default:
			dep, err := fd.deployAddStream.Recv()
			if err != nil {
				log.Printf("[Client] DeployAdd stream ended: %v", err)
				return
			}
			// ADD → Insert
			if err := fd.dbHandler.InsertDeploy(dep); err != nil {
				log.Printf("[MongoDB] InsertDeploy(Add) error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// DeployUpdateRoutine Function
func (fd *Feeder) DeployUpdateRoutine() {
	for fd.Running {
		select {
		default:
			dep, err := fd.deployUpdateStream.Recv()
			if err != nil {
				log.Printf("[Client] DeployUpdate stream ended: %v", err)
				return
			}
			// UPDATE → UpdateDeploy
			if err := fd.dbHandler.UpdateDeploy(dep); err != nil {
				log.Printf("[MongoDB] UpdateDeploy error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// DeployDeleteRoutine Function
func (fd *Feeder) DeployDeleteRoutine() {
	for fd.Running {
		select {
		default:
			dep, err := fd.deployDeleteStream.Recv()
			if err != nil {
				log.Printf("[Client] DeployDelete stream ended: %v", err)
				return
			}
			// DELETE → DeleteDeploy
			if err := fd.dbHandler.DeleteDeploy(dep); err != nil {
				log.Printf("[MongoDB] DeleteDeploy error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// PodAddRoutine Function
func (fd *Feeder) PodAddRoutine() {
	for fd.Running {
		select {
		default:
			pod, err := fd.podAddStream.Recv()
			if err != nil {
				log.Printf("[Client] PodAdd stream ended: %v", err)
				return
			}
			if err := fd.dbHandler.InsertPod(pod); err != nil {
				log.Printf("[MongoDB] InsertPod(Add) error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// PodUpdateRoutine Function
func (fd *Feeder) PodUpdateRoutine() {
	for fd.Running {
		select {
		default:
			pod, err := fd.podUpdateStream.Recv()
			if err != nil {
				log.Printf("[Client] PodUpdate stream ended: %v", err)
				return
			}
			if err := fd.dbHandler.UpdatePod(pod); err != nil {
				log.Printf("[MongoDB] UpdatePod error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// PodDeleteRoutine Function
func (fd *Feeder) PodDeleteRoutine() {
	for fd.Running {
		select {
		default:
			pod, err := fd.podDeleteStream.Recv()
			if err != nil {
				log.Printf("[Client] PodDelete stream ended: %v", err)
				return
			}
			if err := fd.dbHandler.DeletePod(pod); err != nil {
				log.Printf("[MongoDB] DeletePod error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// ServiceAddRoutine Function
func (fd *Feeder) ServiceAddRoutine() {
	for fd.Running {
		select {
		default:
			svc, err := fd.svcAddStream.Recv()
			if err != nil {
				log.Printf("[Client] ServiceAdd stream ended: %v", err)
				return
			}
			if err := fd.dbHandler.InsertService(svc); err != nil {
				log.Printf("[MongoDB] InsertService(Add) error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// ServiceUpdateRoutine Function
func (fd *Feeder) ServiceUpdateRoutine() {
	for fd.Running {
		select {
		default:
			svc, err := fd.svcUpdateStream.Recv()
			if err != nil {
				log.Printf("[Client] ServiceUpdate stream ended: %v", err)
				return
			}
			if err := fd.dbHandler.UpdateService(svc); err != nil {
				log.Printf("[MongoDB] UpdateService error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}

// ServiceDeleteRoutine Function
func (fd *Feeder) ServiceDeleteRoutine() {
	for fd.Running {
		select {
		default:
			svc, err := fd.svcDeleteStream.Recv()
			if err != nil {
				log.Printf("[Client] ServiceDelete stream ended: %v", err)
				return
			}
			if err := fd.dbHandler.DeleteService(svc); err != nil {
				log.Printf("[MongoDB] DeleteService error: %v", err)
			}
		case <-fd.Done:
			return
		}
	}
}
