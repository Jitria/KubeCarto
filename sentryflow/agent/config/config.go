// SPDX-License-Identifier: Apache-2.0

package config

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// AgentConfig structure
type AgentConfig struct {
	CollectorAddr string // Address for Collector gRPC
	CollectorPort string // Port for Collector gRPC

	OperatorAddr string // IP address to use for Operator gRPC
	OperatorPort string // Port to use for Operator gRPC

	ClusterName string // Name of the cluster

	PatchingNamespaces           bool // Enable/Disable patching namespaces with 'istio-injection'
	RestartingPatchedDeployments bool // Enable/Disable restarting deployments after patching

	AggregationPeriod int // Period for aggregating metrics
	CleanUpPeriod     int // Period for cleaning up outdated metrics

	Debug bool // Enable/Disable Agent debug mode
}

// GlobalConfig Global configuration for Agent
var GlobalConfig AgentConfig

// init Function
func init() {
	_ = LoadConfig()
}

// Config const
const (
	CollectorAddr string = "collectorAddr"
	CollectorPort string = "collectorPort"

	OperatorAddr string = "operatorAddr"
	OperatorPort string = "operatorPort"

	ClusterName string = "clusterName"

	PatchingNamespaces           string = "patchingNamespaces"
	RestartingPatchedDeployments string = "restartingPatchedDeployments"

	AggregationPeriod string = "aggregationPeriod"
	CleanUpPeriod     string = "cleanUpPeriod"

	Debug string = "debug"
)

func readCmdLineParams() {
	collectorAddrStr := flag.String(CollectorAddr, "0.0.0.0", "Address for Collector gRPC")
	collectorPortStr := flag.String(CollectorPort, "4317", "Port for Collector gRPC")

	operatorAddrStr := flag.String(OperatorAddr, "sentryflow-operator.sentryflow.svc.cluster.local", "Address for Operator gRPC")
	operatorPortStr := flag.String(OperatorPort, "5317", "Port for Operator gRPC")

	clusterNameStr := flag.String(ClusterName, "", "Name of the Kubernetes cluster")

	patchingNamespacesB := flag.Bool(PatchingNamespaces, false, "Enable patching 'istio-injection' to all namespaces")
	restartingPatchedDeploymentsB := flag.Bool(RestartingPatchedDeployments, false, "Enable restarting the deployments in all patched namespaces")

	aggregationPeriodInt := flag.Int(AggregationPeriod, 1, "Period for aggregating metrics")
	cleanUpPeriodInt := flag.Int(CleanUpPeriod, 5, "Period for cleanning up outdated metrics")

	configDebugB := flag.Bool(Debug, false, "Enable debugging mode")

	var flags []string
	flag.VisitAll(func(f *flag.Flag) {
		kv := fmt.Sprintf("%s:%v", f.Name, f.Value)
		flags = append(flags, kv)
	})
	log.Printf("Arguments [%s]", strings.Join(flags, " "))

	flag.Parse()

	viper.SetDefault(CollectorAddr, *collectorAddrStr)
	viper.SetDefault(CollectorPort, *collectorPortStr)

	viper.SetDefault(OperatorAddr, *operatorAddrStr)
	viper.SetDefault(OperatorPort, *operatorPortStr)

	viper.SetDefault(ClusterName, *clusterNameStr)

	viper.SetDefault(PatchingNamespaces, *patchingNamespacesB)
	viper.SetDefault(RestartingPatchedDeployments, *restartingPatchedDeploymentsB)

	viper.SetDefault(AggregationPeriod, *aggregationPeriodInt)
	viper.SetDefault(CleanUpPeriod, *cleanUpPeriodInt)

	viper.SetDefault(Debug, *configDebugB)
}

// LoadConfig Load configuration
func LoadConfig() error {
	// Read configuration from command line
	readCmdLineParams()

	// Read environment variable, those are upper-cased
	viper.AutomaticEnv()

	GlobalConfig.CollectorAddr = viper.GetString(CollectorAddr)
	GlobalConfig.CollectorPort = viper.GetString(CollectorPort)

	GlobalConfig.OperatorAddr = viper.GetString(OperatorAddr)
	GlobalConfig.OperatorPort = viper.GetString(OperatorPort)

	GlobalConfig.ClusterName = viper.GetString(ClusterName)

	GlobalConfig.PatchingNamespaces = viper.GetBool(PatchingNamespaces)
	GlobalConfig.RestartingPatchedDeployments = viper.GetBool(RestartingPatchedDeployments)

	GlobalConfig.AggregationPeriod = viper.GetInt(AggregationPeriod)
	GlobalConfig.CleanUpPeriod = viper.GetInt(CleanUpPeriod)

	GlobalConfig.Debug = viper.GetBool(Debug)

	log.Printf("Configuration [%+v]", GlobalConfig)

	return nil
}
