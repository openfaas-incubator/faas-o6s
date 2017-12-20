package server

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	clientset "github.com/openfaas-incubator/faas-o6s/pkg/client/clientset/versioned"
	"github.com/openfaas/faas/gateway/requests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeListHandler(namespace string, client clientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		functions := []requests.Function{}

		opts := metav1.ListOptions{}
		res, err := client.O6sV1alpha1().Functions(namespace).List(opts)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			glog.Errorf("Function listing error: %v", err)
			return
		}

		for _, item := range res.Items {
			function := requests.Function{
				Name:     item.Spec.Name,
				Replicas: uint64(*item.Spec.Replicas),
				Image:    item.Spec.Image,
				Labels:   item.Spec.Labels,
			}

			functions = append(functions, function)
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}
