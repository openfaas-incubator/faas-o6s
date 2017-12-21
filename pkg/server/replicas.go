package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	clientset "github.com/openfaas-incubator/faas-o6s/pkg/client/clientset/versioned"
	"github.com/openfaas/faas-provider/types"
	"github.com/openfaas/faas/gateway/requests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeReplicaReader(namespace string, client clientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		functionName := vars["name"]

		opts := metav1.GetOptions{}
		k8sfunc, err := client.O6sV1alpha1().Functions(namespace).Get(functionName, opts)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		result := &requests.Function{
			Replicas:   uint64(*k8sfunc.Spec.Replicas),
			Labels:     k8sfunc.Spec.Labels,
			Name:       k8sfunc.Spec.Name,
			EnvProcess: k8sfunc.Spec.Handler,
			Image:      k8sfunc.Spec.Image,
		}

		res, _ := json.Marshal(result)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	}
}

func makeReplicaHandler(namespace string, client clientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		functionName := vars["name"]

		req := types.ScaleServiceRequest{}
		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			if err := json.Unmarshal(bytesIn, &req); err != nil {
				glog.Errorf("Function %s replica invalid JSON: %v", functionName, err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		opts := metav1.GetOptions{}
		k8sfunc, err := client.O6sV1alpha1().Functions(namespace).Get(functionName, opts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			glog.Errorf("Function %s get error: %v", functionName, err)
			return
		}

		k8sfunc.Spec.Replicas = int32p(int32(req.Replicas))
		_, err = client.O6sV1alpha1().Functions(namespace).Update(k8sfunc)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			glog.Errorf("Function %s update error: %v", functionName, err)
			return
		}

		glog.Infof("Function %v replica updated to %v", functionName, req.Replicas)
		w.WriteHeader(http.StatusAccepted)
	}
}
