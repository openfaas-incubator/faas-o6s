package specutils

import (
	"testing"

	faasv1 "github.com/openfaas-incubator/openfaas-operator/pkg/apis/openfaas/v1alpha2"
)

func Test_MakeReplicas(t *testing.T) {

	nonNilScenarios := []struct {
		name     string
		function *faasv1.Function
		expected int32
	}{
		{"1 replica for the empty function", &faasv1.Function{}, 1},
		{"return original replica value when scaling factor is explicitly 0", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{LabelScalingFactor: "0"}, Replicas: int32p(110)}}, 110},
		{"return original replica value when scaling factor is negative", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{LabelScalingFactor: "-10"}, Replicas: int32p(110)}}, 110},
		{"enforce labels when scaling factor is non-zero", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{LabelScalingFactor: "25", LabelMaxReplicas: "10"}, Replicas: int32p(110)}}, 10},
		{"too small replica returns default min", &faasv1.Function{Spec: faasv1.FunctionSpec{Replicas: int32p(0)}}, 1},
		{"too small replica returns explicit min", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{LabelMinReplicas: "5"}, Replicas: int32p(2)}}, 5},
		{"too large replica returns default max", &faasv1.Function{Spec: faasv1.FunctionSpec{Replicas: int32p(1000)}}, 100},
		{"too large replica returns explicit max", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{LabelMaxReplicas: "20"}, Replicas: int32p(1000)}}, 20},
		{"replica already between explicit min and max is unchanged", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{LabelMinReplicas: "5", LabelMaxReplicas: "10"}, Replicas: int32p(7)}}, 7},
		{"replica is updated to be between explicit min and max", &faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{LabelMinReplicas: "5", LabelMaxReplicas: "10"}, Replicas: int32p(1)}}, 5},
	}

	for _, s := range nonNilScenarios {
		t.Run(s.name, func(t *testing.T) {
			replicas := MakeReplicas(s.function)
			value := *replicas
			if s.expected != value {
				t.Errorf("incorrect replica count: expected %v, got %v", s.expected, value)
			}
		})
	}

	t.Run("allow nil replica value when factor is zero", func(t *testing.T) {
		replicas := MakeReplicas(&faasv1.Function{Spec: faasv1.FunctionSpec{Labels: &map[string]string{LabelScalingFactor: "0"}, Replicas: nil}})
		if replicas != nil {
			t.Errorf("expected nil, got: %v", replicas)
		}
	})
}
func Test_GetMaxReplicas(t *testing.T) {
	scenarios := []struct {
		name     string
		labels   *map[string]string
		expected int32
	}{
		{"nil labels returns default max", nil, int32(defaultMaxReplicas)},
		{"empty labels returns default max", &map[string]string{}, int32(defaultMaxReplicas)},
		{"0 value returns default max", &map[string]string{LabelMaxReplicas: "0"}, int32(defaultMaxReplicas)},
		{"negative value returns default max", &map[string]string{LabelMaxReplicas: "-10"}, int32(defaultMaxReplicas)},
		{"non-negative value returns supplied value", &map[string]string{LabelMaxReplicas: "10"}, int32(10)},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			value := GetMaxReplicaCount(s.labels)
			if s.expected != value {
				t.Errorf("incorrect max replica count: expected %v, got %v", s.expected, value)
			}
		})
	}
}

func Test_GetMinReplicas(t *testing.T) {
	scenarios := []struct {
		name     string
		labels   *map[string]string
		expected int32
	}{
		{"nil labels returns default min", nil, int32(defaultMinReplicas)},
		{"empty labels returns default min", &map[string]string{}, int32(defaultMinReplicas)},
		{"0 value returns default min", &map[string]string{LabelMinReplicas: "0"}, int32(defaultMinReplicas)},
		{"negative value returns default min", &map[string]string{LabelMinReplicas: "-10"}, int32(defaultMinReplicas)},
		{"non-negative value returns supplied value", &map[string]string{LabelMinReplicas: "10"}, int32(10)},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			value := GetMinReplicaCount(s.labels)
			if s.expected != value {
				t.Errorf("incorrect min replica count: expected %v, got %v", s.expected, value)
			}
		})
	}
}

func Test_GetScalingFactor(t *testing.T) {
	scenarios := []struct {
		name     string
		labels   *map[string]string
		expected int
	}{
		{"nil labels returns default factor", nil, defaultScalingFactor},
		{"empty labels returns default factor", &map[string]string{}, defaultScalingFactor},
		{"non-integer returns default factor", &map[string]string{LabelScalingFactor: "test"}, defaultScalingFactor},
		{"0 value returns 0", &map[string]string{LabelScalingFactor: "0"}, 0},
		{"negative value returns 0", &map[string]string{LabelScalingFactor: "-10"}, 0},
		{"non-negative value returns supplied value", &map[string]string{LabelScalingFactor: "10"}, 10},
		{"greater that 100 value returns 100", &map[string]string{LabelScalingFactor: "1000"}, 100},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			value := GetScalingFactor(s.labels)
			if s.expected != value {
				t.Errorf("incorrect replicas scaling factor: expected %v, got %v", s.expected, value)
			}
		})
	}
}

func Test_GetScaleToZero(t *testing.T) {
	scenarios := []struct {
		name     string
		labels   *map[string]string
		expected bool
	}{
		{"nil labels returns default factor", nil, defaultScaleToZero},
		{"empty labels returns default factor", &map[string]string{}, defaultScaleToZero},
		{"non-integer returns default factor", &map[string]string{LabelScaleToZero: "test"}, defaultScaleToZero},
		{"true returns true", &map[string]string{LabelScaleToZero: "true"}, true},
		{"false returns false", &map[string]string{LabelScaleToZero: "false"}, false},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			value := GetScaleToZero(s.labels)
			if s.expected != value {
				t.Errorf("incorrect replicas scaling factor: expected %v, got %v", s.expected, value)
			}
		})
	}
}

func int32p(i int32) *int32 {
	return &i
}
