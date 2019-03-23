package specutils

import (
	"strconv"

	"github.com/golang/glog"
	"github.com/openfaas/faas/gateway/requests"

	v1alpha1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
)

const (
	defaultMinReplicas = 1
	defaultMaxReplicas = 100
	LabelMinReplicas   = "com.openfaas.scale.min"
	LabelMaxReplicas   = "com.openfaas.scale.max"
)

// GetMinReplicaCount parses the function min allowed replicas value form its labels
func GetMinReplicaCount(labels *map[string]string) int32 {
	value := getLabelInt(labels, LabelMinReplicas)
	if value > 0 {
		return int32(value)
	}

	return defaultMinReplicas
}

// GetMaxReplicaCount parses the function max allowed replicas value form its labels
func GetMaxReplicaCount(labels *map[string]string) int32 {
	value := getLabelInt(labels, LabelMaxReplicas)

	if value > 0 {
		return int32(value)
	}

	return defaultMaxReplicas
}

func getLabelInt(labels *map[string]string, name string) int {
	if labels == nil {
		return 0
	}

	lb := *labels
	valueStr, exists := lb[name]
	if !exists {
		return 0
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		glog.Error(err)
		return 0
	}
	return value
}

// GetResources cases the requests.FunctionResources to a k8s v1alpha1.FunctionResources
func GetResources(limits *requests.FunctionResources) *v1alpha1.FunctionResources {
	if limits == nil {
		return nil
	}
	return &v1alpha1.FunctionResources{
		CPU:    limits.CPU,
		Memory: limits.Memory,
	}
}

func int32p(i int32) *int32 {
	return &i
}
