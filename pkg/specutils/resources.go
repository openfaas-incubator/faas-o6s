package specutils

import (
	"github.com/openfaas/faas/gateway/requests"

	v1alpha1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
)

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
