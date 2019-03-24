package specutils

import (
	"testing"

	"github.com/openfaas/faas/gateway/requests"
)

func Test_GetResources(t *testing.T) {
	scenarios := []struct {
		name           string
		request        *requests.FunctionResources
		expectedCPU    string
		expectedMemory string
	}{

		{"nil labels returns nil", nil, "", ""},
		{"CPU only resource request", &requests.FunctionResources{CPU: "1"}, "1", ""},
		{"Memory only resource request", &requests.FunctionResources{Memory: "100"}, "", "100"},
	}
	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			value := GetResources(s.request)
			if s.request == nil {
				if value != nil {
					t.Errorf("expected nil resources, got: %+v", value)
				}
				return
			}

			if value == nil {
				t.Errorf("unexpected nil resources")
			}

			if s.expectedCPU != value.CPU {
				t.Errorf("incorrect cpu value: expected %v, got %v", s.expectedCPU, value.CPU)
			}

			if s.expectedMemory != value.Memory {
				t.Errorf("incorrect memory value: expected %v, got %v", s.expectedMemory, value.Memory)
			}
		})
	}
}
