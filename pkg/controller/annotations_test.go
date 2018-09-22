package controller

import (
	"testing"

	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
)

func Test_makeAnnotations_NoKeys(t *testing.T) {
	annotationVal := `{"name":"","image":"","replicas":null,"handler":"","annotations":null,"labels":null,"environment":null,"constraints":null,"secrets":null,"limits":null,"requests":null,"readOnlyRootFilesystem":false}`

	spec := faasv1.Function{
		Spec: faasv1.FunctionSpec{},
	}

	annotations := makeAnnotations(&spec)

	if _, ok := annotations["prometheus.io.scrape"]; !ok {
		t.Errorf("wanted annotation " + "prometheus.io.scrape" + " to be added")
		t.Fail()
	}
	if val, _ := annotations["prometheus.io.scrape"]; val != "false" {
		t.Errorf("wanted annotation " + "prometheus.io.scrape" + ` to equal "false"`)
		t.Fail()
	}

	if _, ok := annotations[annotationFunctionSpec]; !ok {
		t.Errorf("wanted annotation " + annotationFunctionSpec)
		t.Fail()
	}

	if val, _ := annotations[annotationFunctionSpec]; val != annotationVal {
		t.Errorf("Annotation " + annotationFunctionSpec + "\nwant: '" + annotationVal + "'\ngot: '" + val + "'")
		t.Fail()
	}
}

func Test_makeAnnotations_WithKeyAndValue(t *testing.T) {
	annotationVal := `{"name":"","image":"","replicas":null,"handler":"","annotations":{"key":"value","key2":"value2"},"labels":null,"environment":null,"constraints":null,"secrets":null,"limits":null,"requests":null,"readOnlyRootFilesystem":false}`

	spec := faasv1.Function{
		Spec: faasv1.FunctionSpec{
			Annotations: &map[string]string{
				"key":  "value",
				"key2": "value2",
			},
		},
	}

	annotations := makeAnnotations(&spec)

	if _, ok := annotations["prometheus.io.scrape"]; !ok {
		t.Errorf("wanted annotation " + "prometheus.io.scrape" + " to be added")
		t.Fail()
	}
	if val, _ := annotations["prometheus.io.scrape"]; val != "false" {
		t.Errorf("wanted annotation " + "prometheus.io.scrape" + ` to equal "false"`)
		t.Fail()
	}

	if _, ok := annotations[annotationFunctionSpec]; !ok {
		t.Errorf("wanted annotation " + annotationFunctionSpec)
		t.Fail()
	}

	if val, _ := annotations[annotationFunctionSpec]; val != annotationVal {
		t.Errorf("Annotation " + annotationFunctionSpec + "\nwant: '" + annotationVal + "'\ngot: '" + val + "'")
		t.Fail()
	}
}
