/*
Copyright 2017 OpenFaaS Project

Licensed under the MIT license. See LICENSE file in the project root for full license information.
*/
package v1alpha1

import (
	v1alpha1 "github.com/openfaas-incubator/faas-o6s/pkg/apis/o6sio/v1alpha1"
	"github.com/openfaas-incubator/faas-o6s/pkg/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type O6sV1alpha1Interface interface {
	RESTClient() rest.Interface
	FunctionsGetter
}

// O6sV1alpha1Client is used to interact with features provided by the o6s.io group.
type O6sV1alpha1Client struct {
	restClient rest.Interface
}

func (c *O6sV1alpha1Client) Functions(namespace string) FunctionInterface {
	return newFunctions(c, namespace)
}

// NewForConfig creates a new O6sV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*O6sV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &O6sV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new O6sV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *O6sV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new O6sV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *O6sV1alpha1Client {
	return &O6sV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *O6sV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
