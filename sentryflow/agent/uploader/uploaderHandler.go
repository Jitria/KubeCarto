// SPDX-License-Identifier: Apache-2.0

package uploader

import (
	"fmt"
	"log"
	"sync"

	"Agent/config"

	"github.com/Jitria/SentryFlow/protobuf"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// == //

// UplH global reference for Uploader Handler
var UplH *UplHandler

// init Function
func init() {
	UplH = NewUploaderHandler()
}

// UplHandler Structure
type UplHandler struct {
	grpcClient protobuf.SentryFlowClient

	uploaderAPILogs      chan *protobuf.APILog
	uploaderEnovyMetrics chan *protobuf.EnvoyMetrics

	stopChan chan struct{}
}

// NewUploaderHandler Function
func NewUploaderHandler() *UplHandler {
	ch := &UplHandler{
		uploaderAPILogs:      make(chan *protobuf.APILog),
		uploaderEnovyMetrics: make(chan *protobuf.EnvoyMetrics),

		stopChan: make(chan struct{}),
	}
	return ch
}

// == //

// StartUploader Function
func StartUploader(wg *sync.WaitGroup) bool {
	grpcClient, err := connectToOperator()
	if err != nil {
		log.Printf("[Uploader] Failed to connect to Operator's gRPC server")
		return false
	}
	UplH.grpcClient = grpcClient

	// Export APILogs
	go UplH.uploadAPILogs(wg)
	log.Printf("[Uploader] Exporting API logs through gRPC services")

	// Export EnvoyMetrics
	go UplH.uploadEnvoyMetrics(wg)
	log.Printf("[Uploader] Exporting Envoy metrics through gRPC services")

	return true
}

func connectToOperator() (protobuf.SentryFlowClient, error) {
	operatorAddr := fmt.Sprintf("%s:%s", config.GlobalConfig.OperatorAddr, config.GlobalConfig.OperatorPort)

	conn, err := grpc.NewClient(operatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("[Uploader] Failed to connect to Operator's gRPC server at %s: %v", operatorAddr, err)
		return nil, nil
	}

	client := protobuf.NewSentryFlowClient(conn)

	return client, nil
}

// == //

// StopUploader Function
func StopUploader() bool {
	// One for uploadAPILogs
	UplH.stopChan <- struct{}{}

	// One for uploadEnvoyMetrics
	UplH.stopChan <- struct{}{}

	return true
}

// == //
