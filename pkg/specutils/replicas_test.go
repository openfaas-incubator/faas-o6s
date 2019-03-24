package specutils

import "testing"

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
