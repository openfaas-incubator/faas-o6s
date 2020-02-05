/*
Copyright 2019-2020 OpenFaaS Authors

Licensed under the MIT license. See LICENSE file in the project root for full license information.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	openfaasv1 "github.com/openfaas/openfaas-operator/pkg/apis/openfaas/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeFunctions implements FunctionInterface
type FakeFunctions struct {
	Fake *FakeOpenfaasV1
	ns   string
}

var functionsResource = schema.GroupVersionResource{Group: "openfaas.com", Version: "v1", Resource: "functions"}

var functionsKind = schema.GroupVersionKind{Group: "openfaas.com", Version: "v1", Kind: "Function"}

// Get takes name of the function, and returns the corresponding function object, and an error if there is any.
func (c *FakeFunctions) Get(name string, options v1.GetOptions) (result *openfaasv1.Function, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(functionsResource, c.ns, name), &openfaasv1.Function{})

	if obj == nil {
		return nil, err
	}
	return obj.(*openfaasv1.Function), err
}

// List takes label and field selectors, and returns the list of Functions that match those selectors.
func (c *FakeFunctions) List(opts v1.ListOptions) (result *openfaasv1.FunctionList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(functionsResource, functionsKind, c.ns, opts), &openfaasv1.FunctionList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &openfaasv1.FunctionList{ListMeta: obj.(*openfaasv1.FunctionList).ListMeta}
	for _, item := range obj.(*openfaasv1.FunctionList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested functions.
func (c *FakeFunctions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(functionsResource, c.ns, opts))

}

// Create takes the representation of a function and creates it.  Returns the server's representation of the function, and an error, if there is any.
func (c *FakeFunctions) Create(function *openfaasv1.Function) (result *openfaasv1.Function, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(functionsResource, c.ns, function), &openfaasv1.Function{})

	if obj == nil {
		return nil, err
	}
	return obj.(*openfaasv1.Function), err
}

// Update takes the representation of a function and updates it. Returns the server's representation of the function, and an error, if there is any.
func (c *FakeFunctions) Update(function *openfaasv1.Function) (result *openfaasv1.Function, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(functionsResource, c.ns, function), &openfaasv1.Function{})

	if obj == nil {
		return nil, err
	}
	return obj.(*openfaasv1.Function), err
}

// Delete takes name of the function and deletes it. Returns an error if one occurs.
func (c *FakeFunctions) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(functionsResource, c.ns, name), &openfaasv1.Function{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeFunctions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(functionsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &openfaasv1.FunctionList{})
	return err
}

// Patch applies the patch and returns the patched function.
func (c *FakeFunctions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *openfaasv1.Function, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(functionsResource, c.ns, name, pt, data, subresources...), &openfaasv1.Function{})

	if obj == nil {
		return nil, err
	}
	return obj.(*openfaasv1.Function), err
}
