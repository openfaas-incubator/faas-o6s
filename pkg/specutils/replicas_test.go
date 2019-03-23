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
