package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	v1alpha1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
	clientset "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned"
	"github.com/openfaas-incubator/openfaas-operator/pkg/specutils"
	"github.com/openfaas/faas/gateway/requests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	glog "k8s.io/klog"
)

func makeApplyHandler(namespace string, client clientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Body != nil {
			defer r.Body.Close()
		}

		body, _ := ioutil.ReadAll(r.Body)
		req := requests.CreateFunctionRequest{}
		err := json.Unmarshal(body, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		newFunc := &v1alpha1.Function{
			ObjectMeta: metav1.ObjectMeta{
				Name:      req.Service,
				Namespace: namespace,
			},
			Spec: v1alpha1.FunctionSpec{
				Name:                   req.Service,
				Image:                  req.Image,
				Handler:                req.EnvProcess,
				Labels:                 req.Labels,
				Annotations:            req.Annotations,
				Environment:            &req.EnvVars,
				Constraints:            req.Constraints,
				Secrets:                req.Secrets,
				Replicas:               int32p(specutils.GetMinReplicaCount(req.Labels)),
				Limits:                 specutils.GetResources(req.Limits),
				Requests:               specutils.GetResources(req.Requests),
				ReadOnlyRootFilesystem: req.ReadOnlyRootFilesystem,
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
					w.Write([]byte(err.Error()))
					glog.Errorf("Function %s create error: %v", req.Service, err)
					return
				} else {
					glog.Infof("Function %s created", req.Service)
				}
			} else {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				glog.Errorf("Function %s update error: %v", req.Service, err)
				return
			}
		} else {
			glog.Infof("Function %s updated", req.Service)
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
