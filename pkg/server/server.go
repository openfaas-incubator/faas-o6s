package server

import (
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
	clientset "github.com/openfaas-incubator/faas-o6s/pkg/client/clientset/versioned"
	"github.com/openfaas/faas-provider"
	"github.com/openfaas/faas-provider/types"
)

func Start(client clientset.Interface) {
	functionNamespace := "default"
	if namespace, exists := os.LookupEnv("function_namespace"); exists {
		functionNamespace = namespace
	}

	port := 9090
	if portVal, exists := os.LookupEnv("port"); exists {
		parsedVal, parseErr := strconv.Atoi(portVal)
		if parseErr == nil && parsedVal > 0 {
			port = parsedVal
		}
	}

	readTimeout := 8
	if val, exists := os.LookupEnv("read_timeout"); exists {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal > 0 {
			readTimeout = parsedVal
		}
	}

	writeTimeout := 8
	if val, exists := os.LookupEnv("write_timeout"); exists {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal > 0 {
			writeTimeout = parsedVal
		}
	}

	bootstrapHandlers := types.FaaSHandlers{
		FunctionProxy:  makeProxy(functionNamespace),
		DeleteHandler:  makeDeleteHandler(functionNamespace, client),
		DeployHandler:  makeApplyHandler(functionNamespace, client),
		FunctionReader: makeListHandler(functionNamespace, client),
		ReplicaReader:  makeReplicaReader(functionNamespace, client),
		ReplicaUpdater: makeReplicaHandler(functionNamespace, client),
		UpdateHandler:  makeApplyHandler(functionNamespace, client),
	}

	bootstrapConfig := types.FaaSConfig{
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		TCPPort:      &port,
	}

	glog.Infof("Using namespace %v", functionNamespace)
	glog.Infof("Starting HTTP server on port %v", port)
	bootstrap.Serve(&bootstrapHandlers, &bootstrapConfig)
}
