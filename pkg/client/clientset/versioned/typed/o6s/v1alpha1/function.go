/*
Copyright 2017 OpenFaaS Project

Licensed under the MIT license. See LICENSE file in the project root for full license information.
*/
package v1alpha1

import (
	v1alpha1 "github.com/openfaas-incubator/faas-o6s/pkg/apis/o6sio/v1alpha1"
	scheme "github.com/openfaas-incubator/faas-o6s/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// FunctionsGetter has a method to return a FunctionInterface.
// A group's client should implement this interface.
type FunctionsGetter interface {
	Functions(namespace string) FunctionInterface
}

// FunctionInterface has methods to work with Function resources.
type FunctionInterface interface {
	Create(*v1alpha1.Function) (*v1alpha1.Function, error)
	Update(*v1alpha1.Function) (*v1alpha1.Function, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Function, error)
	List(opts v1.ListOptions) (*v1alpha1.FunctionList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Function, err error)
	FunctionExpansion
}

// functions implements FunctionInterface
type functions struct {
	client rest.Interface
	ns     string
}

// newFunctions returns a Functions
func newFunctions(c *O6sV1alpha1Client, namespace string) *functions {
	return &functions{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the function, and returns the corresponding function object, and an error if there is any.
func (c *functions) Get(name string, options v1.GetOptions) (result *v1alpha1.Function, err error) {
	result = &v1alpha1.Function{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("functions").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Functions that match those selectors.
func (c *functions) List(opts v1.ListOptions) (result *v1alpha1.FunctionList, err error) {
	result = &v1alpha1.FunctionList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("functions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested functions.
func (c *functions) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("functions").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a function and creates it.  Returns the server's representation of the function, and an error, if there is any.
func (c *functions) Create(function *v1alpha1.Function) (result *v1alpha1.Function, err error) {
	result = &v1alpha1.Function{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("functions").
		Body(function).
		Do().
		Into(result)
	return
}

// Update takes the representation of a function and updates it. Returns the server's representation of the function, and an error, if there is any.
func (c *functions) Update(function *v1alpha1.Function) (result *v1alpha1.Function, err error) {
	result = &v1alpha1.Function{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("functions").
		Name(function.Name).
		Body(function).
		Do().
		Into(result)
	return
}

// Delete takes name of the function and deletes it. Returns an error if one occurs.
func (c *functions) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("functions").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *functions) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("functions").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched function.
func (c *functions) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Function, err error) {
	result = &v1alpha1.Function{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("functions").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
