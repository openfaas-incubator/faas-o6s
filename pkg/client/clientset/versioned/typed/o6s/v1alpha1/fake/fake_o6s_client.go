/*
Copyright 2017 OpenFaaS Project

Licensed under the MIT license. See LICENSE file in the project root for full license information.
*/
package fake

import (
	v1alpha1 "github.com/openfaas-incubator/faas-o6s/pkg/client/clientset/versioned/typed/o6s/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeO6sV1alpha1 struct {
	*testing.Fake
}

func (c *FakeO6sV1alpha1) Functions(namespace string) v1alpha1.FunctionInterface {
	return &FakeFunctions{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeO6sV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
