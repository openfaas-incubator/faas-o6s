package server

import (
	"encoding/json"
	"github.com/openfaas/faas/gateway/requests"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	glog "k8s.io/klog"
	"net/http"
)

// makeSecretHandler provides the secrets CRUD endpoint
func makeSecretHandler(namespace string, kube kubernetes.Interface, secretLister corev1listers.SecretNamespaceLister) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

		switch r.Method {
		case http.MethodGet:
			ls := map[string]string{
				"owner": "openfaas",
			}
			res, err := secretLister.List(labels.SelectorFromSet(ls))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				glog.Errorf("Secrets query error: %v", err)
				return
			}
			secrets := []requests.Secret{}
			for _, item := range res {
				secret := requests.Secret{
					Name: item.Name,
				}
				secrets = append(secrets, secret)
			}
			secretsBytes, err := json.Marshal(secrets)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				glog.Errorf("Secrets json marshal error: %v", err)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(secretsBytes)
		case http.MethodPost:
			secret, err := getSecret(namespace, r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				glog.Errorf("Secret unmarshal error: %v", err)
				return
			}
			_, err = kube.CoreV1().Secrets(namespace).Create(secret)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				glog.Errorf("Secret create error: %v", err)
				return
			}
			glog.Infof("Secret %s created", secret.GetName())
			w.WriteHeader(http.StatusAccepted)
		case http.MethodPut:
			newSecret, err := getSecret(namespace, r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				glog.Errorf("Secret unmarshal error: %v", err)
				return
			}
			secret, err := secretLister.Get(newSecret.GetName())
			if errors.IsNotFound(err) {
				w.WriteHeader(http.StatusNotFound)
				glog.Warningf("Secret update error: %s not found", newSecret.GetName())
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				glog.Errorf("Secret query error: %v", err)
				return
			}
			secret.StringData = newSecret.StringData
			_, err = kube.CoreV1().Secrets(namespace).Update(secret)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				glog.Errorf("Secret update error: %v", err)
				return
			}
			glog.Infof("Secret %s updated", secret.GetName())
			w.WriteHeader(http.StatusAccepted)
		case http.MethodDelete:
			secret, err := getSecret(namespace, r.Body)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				glog.Errorf("Secret unmarshal error: %v", err)
				return
			}
			opts := &metav1.DeleteOptions{}
			err = kube.CoreV1().Secrets(namespace).Delete(secret.GetName(), opts)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				glog.Errorf("Secret %s delete error: %v", secret.GetName(), err)
				return
			}
			glog.Infof("Secret %s deleted", secret.GetName())
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}

func getSecret(namespace string, r io.Reader) (*corev1.Secret, error) {
	body, _ := ioutil.ReadAll(r)
	req := requests.Secret{}
	err := json.Unmarshal(body, &req)
	if err != nil {
		return nil, err
	}
	secret := &corev1.Secret{
		Type: corev1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: namespace,
			Labels: map[string]string{
				"owner": "openfaas",
			},
		},
		StringData: map[string]string{
			req.Name: req.Value,
		},
	}

	return secret, nil
}
