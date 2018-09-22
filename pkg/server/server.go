package server

import (
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"

	"github.com/golang/glog"
	clientset "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned"
	"github.com/openfaas/faas-provider"
	"github.com/openfaas/faas-provider/types"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/kubernetes"
)

// TODO: Move to config pattern used else-where across project

const defaultHTTPPort = 8081
const defaultReadTimeout = 8
const defaultWriteTimeout = 8

// Start starts HTTP Server for API
func Start(client clientset.Interface, kube kubernetes.Interface) {
	functionNamespace := "openfaas-fn"
	if namespace, exists := os.LookupEnv("function_namespace"); exists {
		functionNamespace = namespace
	}

	port := defaultHTTPPort
	if portVal, exists := os.LookupEnv("port"); exists {
		parsedVal, parseErr := strconv.Atoi(portVal)
		if parseErr == nil && parsedVal > 0 {
			port = parsedVal
		}
	}

	readTimeout := defaultReadTimeout
	if val, exists := os.LookupEnv("read_timeout"); exists {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal > 0 {
			readTimeout = parsedVal
		}
	}

	writeTimeout := defaultWriteTimeout
	if val, exists := os.LookupEnv("write_timeout"); exists {
		parsedVal, parseErr := strconv.Atoi(val)
		if parseErr == nil && parsedVal > 0 {
			writeTimeout = parsedVal
		}
	}

	pprof := "false"
	if val, exists := os.LookupEnv("pprof"); exists {
		pprof = val
	}

	bootstrapHandlers := types.FaaSHandlers{
		FunctionProxy:  makeProxy(functionNamespace, time.Duration(readTimeout)*time.Second),
		DeleteHandler:  makeDeleteHandler(functionNamespace, client),
		DeployHandler:  makeApplyHandler(functionNamespace, client),
		FunctionReader: makeListHandler(functionNamespace, client, kube),
		ReplicaReader:  makeReplicaReader(functionNamespace, client, kube),
		ReplicaUpdater: makeReplicaHandler(functionNamespace, client),
		UpdateHandler:  makeApplyHandler(functionNamespace, client),
		Health:         makeHealthHandler(),
		InfoHandler:    makeInfoHandler(),
	}

	bootstrapConfig := types.FaaSConfig{
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		TCPPort:      &port,
		EnableHealth: true,
	}

	if pprof == "true" {
		bootstrap.Router().PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	}

	bootstrap.Router().Path("/metrics").Handler(promhttp.Handler())

	glog.Infof("Using namespace '%s'", functionNamespace)
	glog.Infof("Starting HTTP server on port %v", port)
	bootstrap.Serve(&bootstrapHandlers, &bootstrapConfig)
}
