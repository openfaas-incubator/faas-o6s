package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"strconv"
	"strings"

	"github.com/golang/glog"
	v1alpha1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
	clientset "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned"
	"github.com/openfaas/faas/gateway/requests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeApplyHandler(namespace string, client clientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		req := requests.CreateFunctionRequest{}
		err := json.Unmarshal(body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		newFunc := &v1alpha1.Function{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Service,
				Namespace: namespace,
			},
			Spec: v1alpha1.FunctionSpec{
				Name:        req.Service,
				Image:       req.Image,
				Handler:     req.EnvProcess,
				Labels:      req.Labels,
				Environment: &req.EnvVars,
				Constraints: req.Constraints,
				Secrets:     req.Secrets,
				Replicas:    getMinReplicaCount(req.Labels),
				Limits:      getResources(req.Limits),
				Requests:    getResources(req.Requests),
			},
		}

		opts := metav1.GetOptions{}
		oldFunc, _ := client.OpenfaasV1alpha2().Functions(namespace).Get(req.Service, opts)
		if oldFunc != nil {
			newFunc.ResourceVersion = oldFunc.ResourceVersion
		}
		_, err = client.OpenfaasV1alpha2().Functions(namespace).Update(newFunc)
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "not found") {
				_, err = client.OpenfaasV1alpha2().Functions(namespace).Create(newFunc)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					glog.Errorf("Function %s create error: %v", req.Service, err)
					return
				} else {
					glog.Infof("Function %s created", req.Service)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				glog.Errorf("Function %s update error: %v", req.Service, err)
				return
			}
		} else {
			glog.Infof("Function %s updated", req.Service)
		}

		w.WriteHeader(http.StatusAccepted)
	}
}

func getMinReplicaCount(labels *map[string]string) *int32 {
	if labels != nil {
		lb := *labels
		if value, exists := lb["com.openfaas.scale.min"]; exists {
			minReplicas, err := strconv.Atoi(value)
			if err == nil && minReplicas > 0 {
				return int32p(int32(minReplicas))
			} else {
				glog.Error(err)
			}
		}
	}

	return int32p(1)
}

func getResources(limits *requests.FunctionResources) *v1alpha1.FunctionResources {
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
