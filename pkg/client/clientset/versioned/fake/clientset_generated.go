/*
Copyright 2018 OpenFaaS Authors

Licensed under the MIT license. See LICENSE file in the project root for full license information.
*/
package fake

import (
	clientset "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned"
	openfaasv1alpha2 "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned/typed/openfaas/v1alpha2"
	fakeopenfaasv1alpha2 "github.com/openfaas-incubator/openfaas-operator/pkg/client/clientset/versioned/typed/openfaas/v1alpha2/fake"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/discovery"
	fakediscovery "k8s.io/client-go/discovery/fake"
	"k8s.io/client-go/testing"
)

// NewSimpleClientset returns a clientset that will respond with the provided objects.
// It's backed by a very simple object tracker that processes creates, updates and deletions as-is,
// without applying any validations and/or defaults. It shouldn't be considered a replacement
// for a real clientset and is mostly useful in simple unit tests.
func NewSimpleClientset(objects ...runtime.Object) *Clientset {
	o := testing.NewObjectTracker(scheme, codecs.UniversalDecoder())
	for _, obj := range objects {
		if err := o.Add(obj); err != nil {
			panic(err)
		}
	}

	fakePtr := testing.Fake{}
	fakePtr.AddReactor("*", "*", testing.ObjectReaction(o))
	fakePtr.AddWatchReactor("*", testing.DefaultWatchReactor(watch.NewFake(), nil))

	return &Clientset{fakePtr, &fakediscovery.FakeDiscovery{Fake: &fakePtr}}
}

// Clientset implements clientset.Interface. Meant to be embedded into a
// struct to get a default implementation. This makes faking out just the method
// you want to test easier.
type Clientset struct {
	testing.Fake
	discovery *fakediscovery.FakeDiscovery
}

func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	return c.discovery
}

var _ clientset.Interface = &Clientset{}

// OpenfaasV1alpha2 retrieves the OpenfaasV1alpha2Client
func (c *Clientset) OpenfaasV1alpha2() openfaasv1alpha2.OpenfaasV1alpha2Interface {
	return &fakeopenfaasv1alpha2.FakeOpenfaasV1alpha2{Fake: &c.Fake}
}

// Openfaas retrieves the OpenfaasV1alpha2Client
func (c *Clientset) Openfaas() openfaasv1alpha2.OpenfaasV1alpha2Interface {
	return &fakeopenfaasv1alpha2.FakeOpenfaasV1alpha2{Fake: &c.Fake}
}
