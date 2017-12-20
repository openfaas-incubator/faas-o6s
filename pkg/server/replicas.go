package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	clientset "github.com/openfaas-incubator/faas-o6s/pkg/client/clientset/versioned"
	"github.com/openfaas/faas/gateway/requests"
)

func makeReplicaReader(namespace string, client clientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		functions := []requests.Function{}
		vars := mux.Vars(r)
		functionName := vars["name"]

		var found *requests.Function
		for _, function := range functions {
			if function.Name == functionName {
				found = &function
				break
			}
		}
		if found == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		functionBytes, _ := json.Marshal(found)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}

func makeReplicaHandler(namespace string, client clientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		functionName := vars["name"]

		req := ScaleServiceRequest{}
		if r.Body != nil {
			defer r.Body.Close()
			bytesIn, _ := ioutil.ReadAll(r.Body)
			marshalErr := json.Unmarshal(bytesIn, &req)
			if marshalErr != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}

		glog.Infof("Function %v replica updated to %v", functionName, req.Replicas)
	}
}

type ScaleServiceRequest struct {
	ServiceName string `json:"serviceName"`
	Replicas    uint64 `json:"replicas"`
}
