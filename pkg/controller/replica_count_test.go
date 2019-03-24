package controller

import (
	"testing"

	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
	"github.com/openfaas-incubator/openfaas-operator/pkg/specutils"
)

func Test_newDeployment_replica_counts(t *testing.T) {

	scenarios := []struct {
		name     string
		function *faasv1.Function
		expected int32
	}{
		{"1 replica for the empty function", &faasv1.Function{}, 1},
		{"return original replica value when scaling factor is explicitly 0", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{specutils.LabelScalingFactor: "0"}, Replicas: int32p(110)}}, 110},
		{"return original replica value when scaling factor is negative", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{specutils.LabelScalingFactor: "-10"}, Replicas: int32p(110)}}, 110},
		{"enforce labels when scaling factor is non-zero", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{specutils.LabelScalingFactor: "25", specutils.LabelMaxReplicas: "10"}, Replicas: int32p(110)}}, 10},
		{"too small replica returns default min", &faasv1.Function{Spec: faasv1.FunctionSpec{Replicas: int32p(0)}}, 1},
		{"too small replica returns explicit min", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{specutils.LabelMinReplicas: "5"}, Replicas: int32p(2)}}, 5},
		{"too large replica returns default max", &faasv1.Function{Spec: faasv1.FunctionSpec{Replicas: int32p(1000)}}, 100},
		{"too large replica returns explicit max", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{specutils.LabelMaxReplicas: "20"}, Replicas: int32p(1000)}}, 20},
		{"replica already between explicit min and max is unchanged", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{specutils.LabelMinReplicas: "5", specutils.LabelMaxReplicas: "10"}, Replicas: int32p(7)}}, 7},
		{"replica is updated to be between explicit min and max", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{specutils.LabelMinReplicas: "5", specutils.LabelMaxReplicas: "10"}, Replicas: int32p(1)}}, 5},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			deploy := newDeployment(s.function, nil, "")
			value := *deploy.Spec.Replicas
			if s.expected != value {
				t.Errorf("incorrect replica count: expected %v, got %v", s.expected, value)
			}
		})
	}
}
