package controller

import (
	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
	"github.com/openfaas/faas-netes/k8s"
	"github.com/openfaas/faas/gateway/requests"
	appsv1 "k8s.io/api/apps/v1beta2"
	"k8s.io/client-go/kubernetes"
)

// FunctionFactory wraps faas-netes factory
type FunctionFactory struct {
	Factory k8s.FunctionFactory
}

func NewFunctionFactory(clientset kubernetes.Interface, config k8s.DeploymentConfig) FunctionFactory {
	return FunctionFactory{
		k8s.FunctionFactory{
			Client: clientset,
			Config: config,
		},
	}
}

func functionToFunctionRequest(in *faasv1.Function) requests.CreateFunctionRequest {
	env := make(map[string]string)
	if in.Spec.Environment != nil {
		env = *in.Spec.Environment
	}
	lim, req := functionToFunctionResources(in)
	return requests.CreateFunctionRequest{
		Annotations:            in.Spec.Annotations,
		Service:                in.Name,
		Labels:                 &in.Labels,
		Constraints:            in.Spec.Constraints,
		EnvProcess:             in.Spec.Handler,
		EnvVars:                env,
		Image:                  in.Spec.Image,
		Limits:                 lim,
		Requests:               req,
		ReadOnlyRootFilesystem: in.Spec.ReadOnlyRootFilesystem,
	}
}

func functionToFunctionResources(in *faasv1.Function) (l *requests.FunctionResources, r *requests.FunctionResources) {
	if in.Spec.Limits != nil {
		l = &requests.FunctionResources{
			Memory: in.Spec.Limits.Memory,
			CPU:    in.Spec.Limits.CPU,
		}
	}
	if in.Spec.Requests != nil {
		r = &requests.FunctionResources{
			Memory: in.Spec.Requests.Memory,
			CPU:    in.Spec.Requests.CPU,
		}
	}
	return
}

func (f *FunctionFactory) MakeProbes(function *faasv1.Function) (*k8s.FunctionProbes, error) {
	req := functionToFunctionRequest(function)
	return f.Factory.MakeProbes(req)
}

func (f *FunctionFactory) ConfigureReadOnlyRootFilesystem(function *faasv1.Function, deployment *appsv1.Deployment) {
	req := functionToFunctionRequest(function)
	f.Factory.ConfigureReadOnlyRootFilesystem(req, deployment)
}

func (f *FunctionFactory) ConfigureContainerUserID(deployment *appsv1.Deployment) {
	f.Factory.ConfigureContainerUserID(deployment)
}
