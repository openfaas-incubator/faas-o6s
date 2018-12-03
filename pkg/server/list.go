package server

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
	clientset "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned"
	"github.com/openfaas/faas/gateway/requests"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/listers/apps/v1beta2"
)

func makeListHandler(namespace string, client clientset.Interface, kube kubernetes.Interface, deploymentLister v1beta2.DeploymentNamespaceLister) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

		functions := []requests.Function{}

		opts := metav1.ListOptions{}
		res, err := client.OpenfaasV1alpha2().Functions(namespace).List(opts)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			glog.Errorf("Function listing error: %v", err)
			return
		}

		for _, item := range res.Items {

			desiredReplicas, availableReplicas, err := getReplicas(item.Spec.Name, namespace, deploymentLister)
			if err != nil {
				glog.Warningf("Function listing getReplicas error: %v", err)
			}

			function := requests.Function{
				Name:              item.Spec.Name,
				Replicas:          desiredReplicas,
				AvailableReplicas: availableReplicas,
				Image:             item.Spec.Image,
				Labels:            item.Spec.Labels,
				Annotations:       item.Spec.Annotations,
			}

			functions = append(functions, function)
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)

	}
}
