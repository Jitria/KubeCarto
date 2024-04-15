// SPDX-License-Identifier: Apache-2.0

package exporter

import (
	"errors"
	"fmt"
	"net"
	"sort"
	"strings"
	"sync"
	"time"

	cfg "github.com/5GSEC/SentryFlow/config"
	"github.com/5GSEC/SentryFlow/protobuf"
	"github.com/5GSEC/SentryFlow/types"

	"github.com/emicklei/go-restful/v3/log"
	"google.golang.org/grpc"
)

// Exp global reference for Exporter Handler
var Exp *Handler

// init Function
func init() {
	Exp = NewExporterHandler()
}

// Handler structure
type Handler struct {
	baseExecutionID    uint64
	currentLogCount    uint64
	agTime             int
	exTime             int
	stopChan           chan struct{}
	lock               sync.Mutex
	exporters          []*Inform
	apiMetricExporters []*apiMetricStreamInform
	metricExporters    []*metricStreamInform
	exporterLock       sync.Mutex
	exporterLogs       chan *protobuf.APILog
	exporterAPIMetrics chan *protobuf.APIMetric
	exporterMetrics    chan *protobuf.EnvoyMetric

	listener   net.Listener
	gRPCServer *grpc.Server
}

// Inform structure
type Inform struct {
	stream    protobuf.SentryFlow_GetLogServer
	error     chan error
	Hostname  string
	IPAddress string
}

// apiMetricStreamInform structure
type apiMetricStreamInform struct {
	apiMetricStream protobuf.SentryFlow_GetAPIMetricsServer
	error           chan error
	Hostname        string
	IPAddress       string
}

// metricStreamInform structure
type metricStreamInform struct {
	metricStream protobuf.SentryFlow_GetEnvoyMetricsServer
	error        chan error
	Hostname     string
	IPAddress    string
}

// NewExporterHandler Function
func NewExporterHandler() *Handler {
	exp := &Handler{
		baseExecutionID:    uint64(time.Now().UnixMicro()),
		currentLogCount:    0,
		agTime:             cfg.GlobalCfg.MetricsDBAggregationTime,
		exTime:             cfg.GlobalCfg.APIMetricsSendTime,
		exporters:          make([]*Inform, 0),
		stopChan:           make(chan struct{}),
		lock:               sync.Mutex{},
		exporterLock:       sync.Mutex{},
		exporterLogs:       make(chan *protobuf.APILog),
		exporterAPIMetrics: make(chan *protobuf.APIMetric),
		exporterMetrics:    make(chan *protobuf.EnvoyMetric),
	}

	return exp
}

// InsertAccessLog Function
func InsertAccessLog(al *protobuf.APILog) {
	// Avoid race condition for currentLogCount, otherwise we might have duplicate IDs
	Exp.lock.Lock()
	al.Id = Exp.baseExecutionID + Exp.currentLogCount
	Exp.currentLogCount++
	Exp.lock.Unlock()

	go saveAccessLog(al) // go routine??
	Exp.exporterLogs <- al
}

func saveAccessLog(al *protobuf.APILog) {
	curLabels := al.SrcLabel

	var labelString []string

	for key, value := range curLabels {
		labelString = append(labelString, fmt.Sprintf("%s:%s", key, value))
	}

	sort.Strings(labelString)

	curData := types.DbAccessLogType{
		Labels:    strings.Join(labelString, " "),
		Namespace: al.SrcNamespace,
		AccessLog: al,
	}

	err := MDB.AccessLogInsert(curData)
	if err != nil {
		log.Printf("unable to insert AccessLog")
		return
	}
}

// InsertEnvoyMetric Function
func InsertEnvoyMetric(em *protobuf.EnvoyMetric) {
	Exp.exporterMetrics <- em
}

// InitExporterServer Function
func (exp *Handler) InitExporterServer() error {
	listenAddr := fmt.Sprintf("%s:%s", cfg.GlobalCfg.CustomExportListenAddr, cfg.GlobalCfg.CustomExportListenPort)

	// Start listening
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		msg := fmt.Sprintf("unable to listen at %s: %v", listenAddr, err)
		return errors.New(msg)
	}

	// Create gRPC server
	server := grpc.NewServer()
	protobuf.RegisterSentryFlowServer(server, exs)

	exp.listener = lis
	exp.gRPCServer = server

	log.Printf("[Exporter] Exporter listening at %s", listenAddr)
	return nil
}

// StartExporterServer Function
func (exp *Handler) StartExporterServer(wg *sync.WaitGroup) error {
	log.Printf("[Exporter] Starting exporter server")
	var err error
	err = nil

	go exp.exportRoutine(wg)

	go func() {
		wg.Add(1)
		// Serve is blocking function
		err = exp.gRPCServer.Serve(exp.listener)
		if err != nil {
			wg.Done()
			return
		}

		wg.Done()
	}()

	return err
}

// exportRoutine Function
func (exp *Handler) exportRoutine(wg *sync.WaitGroup) {
	wg.Add(1)
	log.Printf("[Exporter] Starting export routine")

routineLoop:
	for {
		select {
		// @todo add more channels for this
		case al, ok := <-exp.exporterLogs:
			if !ok {
				log.Printf("[Exporter] Log exporter channel closed")
				break routineLoop
			}

			err := exp.sendLogs(al)
			if err != nil {
				log.Printf("[Exporter] Log exporting failed %v:", err)
			}

		case em, ok := <-exp.exporterMetrics:
			if !ok {
				log.Printf("[Exporter] EnvoyMetric exporter channel closed")
				break routineLoop
			}

			err := exp.sendMetrics(em)
			if err != nil {
				log.Printf("[Exporter] Metric exporting failed %v:", err)
			}

		case am, ok := <-exp.exporterAPIMetrics:
			if !ok {
				log.Printf("[Exporter] APIMetric exporter channel closed")
				break routineLoop
			}
			err := exp.sendAPIMetrics(am)
			if err != nil {
				log.Printf("[Exporter] APIMetric exporting failed %v:", err)
			}

		case <-exp.stopChan:
			break routineLoop
		}
	}

	defer wg.Done()
	return
}

// sendLogs Function
func (exp *Handler) sendLogs(l *protobuf.APILog) error {
	exp.exporterLock.Lock()
	defer exp.exporterLock.Unlock()

	// iterate and send logs
	failed := 0
	total := len(exp.exporters)
	for _, exporter := range exp.exporters {
		curRetry := 0

		// @todo: make max retry count per logs using config
		// @todo: make max retry count per single exporter before removing the exporter using config
		var err error
		for curRetry < 3 {
			err = exporter.stream.Send(l)
			if err != nil {
				log.Printf("[Exporter] Unable to send log to %s(%s) retry=%d/%d: %v",
					exporter.Hostname, exporter.IPAddress, curRetry, 3, err)
				curRetry++
			} else {
				break
			}
		}

		// Count failed
		if err != nil {
			failed++
		}
	}

	// notify failed count
	if failed != 0 {
		msg := fmt.Sprintf("unable to send logs properly %d/%d failed", failed, total)
		return errors.New(msg)
	}

	return nil
}

// sendMetrics Function
func (exp *Handler) sendMetrics(l *protobuf.EnvoyMetric) error {
	exp.exporterLock.Lock()
	defer exp.exporterLock.Unlock()

	// iterate and send logs
	failed := 0
	total := len(exp.metricExporters)
	for _, exporter := range exp.metricExporters {
		curRetry := 0

		// @todo: make max retry count per logs using config
		// @todo: make max retry count per single exporter before removing the exporter using config
		var err error
		for curRetry < 3 {
			err = exporter.metricStream.Send(l)
			if err != nil {
				log.Printf("[Exporter] Unable to send metric to %s(%s) retry=%d/%d: %v",
					exporter.Hostname, exporter.IPAddress, curRetry, 3, err)
				curRetry++
			} else {
				break
			}
		}

		// Count failed
		if err != nil {
			failed++
		}
	}

	// notify failed count
	if failed != 0 {
		msg := fmt.Sprintf("unable to send metrics properly %d/%d failed", failed, total)
		return errors.New(msg)
	}

	return nil
}

// sendAPIMetrics Function
func (exp *Handler) sendAPIMetrics(l *protobuf.APIMetric) error {
	exp.exporterLock.Lock()
	defer exp.exporterLock.Unlock()

	// iterate and send logs
	failed := 0
	total := len(exp.apiMetricExporters)
	for _, exporter := range exp.apiMetricExporters {
		curRetry := 0

		// @todo: make max retry count per logs using config
		// @todo: make max retry count per single exporter before removing the exporter using config
		var err error
		for curRetry < 3 {
			err = exporter.apiMetricStream.Send(l)
			if err != nil {
				log.Printf("[Exporter] Unable to send metric to %s(%s) retry=%d/%d: %v",
					exporter.Hostname, exporter.IPAddress, curRetry, 3, err)
				curRetry++
			} else {
				break
			}
		}

		// Count failed
		if err != nil {
			failed++
		}
	}

	// notify failed count
	if failed != 0 {
		msg := fmt.Sprintf("unable to send metrics properly %d/%d failed", failed, total)
		return errors.New(msg)
	}

	return nil
}

// APIMetricsExportRoutine function
func (exp *Handler) APIMetricsExportRoutine() {

}

// aggregationTimeTickerRoutine Function
func aggregationTimeTickerRoutine() error {
	aggregationTicker := time.NewTicker(time.Duration(Exp.agTime) * time.Second)

	defer aggregationTicker.Stop()

	for {
		select {
		case <-aggregationTicker.C:
			als, err := MDB.AggregatedAccessLogSelect()
			if err != nil {
				log.Printf("[Exporter] AccessLog Aggregation %v", err)
				return err
			}

			for _, val := range als {
				// export part
				curAPIs := []string{}
				for _, APILog := range val {
					curAPIs = append(curAPIs, APILog.Path)
				}
				InsertAPILog(curAPIs)
			}
		}
	}
}

// exportTimeTickerRoutine Function
func exportTimeTickerRoutine() error {
	apiMetricTicker := time.NewTicker(time.Duration(Exp.exTime) * time.Second)

	defer apiMetricTicker.Stop()

	for {
		select {
		case <-apiMetricTicker.C:
			curAPIMetrics, err := MDB.GetAllMetrics()

			if err != nil {
				log.Printf("[Exporter] APIMetric TimeTicker channel closed")
				return err
			}

			if len(curAPIMetrics) > 0 {
				curAPIMetric := &protobuf.APIMetric{
					PerAPICounts: curAPIMetrics,
				}
				Exp.exporterAPIMetrics <- curAPIMetric
			}
		}
	}
}

// StopExporterServer Function
func (exp *Handler) StopExporterServer() {
	// Gracefully stop all client connections
	exp.stopChan <- struct{}{}

	// Gracefully stop gRPC Server
	exp.gRPCServer.GracefulStop()

	log.Printf("[Exporter] Stopped exporter server")
}
