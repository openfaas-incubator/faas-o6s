package controller

import (
	"testing"

	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
)

func Test_makeAnnotationsDoesNotModifyOriginalSpec(t *testing.T) {
	specAnnotations := map[string]string{
		"test.foo": "bar",
	}
	function := &faasv1.Function{
		Spec: faasv1.FunctionSpec{
			Name:        "testfunc",
			Annotations: &specAnnotations,
		},
	}

	expectedAnnotations := map[string]string{
		"prometheus.io.scrape": "false",
		"test.foo":             "bar",
		annotationFunctionSpec: `{"name":"testfunc","image":"","replicas":null,"handler":"","annotations":{"test.foo":"bar"},"labels":null,"environment":null,"constraints":null,"secrets":null,"limits":null,"requests":null,"readOnlyRootFilesystem":false}`,
	}

	makeAnnotations(function)
	annotations := makeAnnotations(function)

	if len(specAnnotations) != 1 {
		t.Errorf("length of original spec annotations has changed, expected 1, got %d", len(specAnnotations))
	}

	if specAnnotations["test.foo"] != "bar" {
		t.Errorf("original spec annotation has changed")
	}

	for name, expectedValue := range expectedAnnotations {
		actualValue := annotations[name]
		if actualValue != expectedValue {
			t.Fatalf("incorrect annotation for '%s': \nexpected %s,\ngot %s", name, expectedValue, actualValue)
		}
	}
}
