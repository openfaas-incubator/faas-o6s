package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	clientset "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned"
	"github.com/openfaas/faas/gateway/requests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func makeDeleteHandler(namespace string, client clientset.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		request := requests.DeleteFunctionRequest{}
		err := json.Unmarshal(body, &request)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(request.FunctionName) == 0 {
			w.WriteHeader(http.StatusBadRequest)
		}

		opts := &metav1.DeleteOptions{}
		err = client.O6sV1alpha1().Functions(namespace).Delete(request.FunctionName, opts)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			glog.Errorf("Function %s delete error: %v", request.FunctionName, err)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
