package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Function describes an OpenFaaS function
type Function struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FunctionSpec   `json:"spec"`
	Status FunctionStatus `json:"status"`
}

// FunctionSpec is the spec for a Function resource
type FunctionSpec struct {
	Name                   string             `json:"name"`
	Image                  string             `json:"image"`
	Replicas               *int32             `json:"replicas"`
	Handler                string             `json:"handler"`
	Labels                 *map[string]string `json:"labels"`
	Environment            *map[string]string `json:"environment"`
	Constraints            []string           `json:"constraints"`
	Secrets                []string           `json:"secrets"`
	Limits                 *FunctionResources `json:"limits"`
	Requests               *FunctionResources `json:"requests"`
	ReadOnlyRootFilesystem bool               `json:"readOnlyRootFilesystem"`
}

// FunctionResources is used to set CPU and memory limits and requests
type FunctionResources struct {
	Memory string `json:"memory"`
	CPU    string `json:"cpu"`
}

// FunctionStatus is the status for a Function resource
type FunctionStatus struct {
	AvailableReplicas int32 `json:"availableReplicas"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FunctionList is a list of Function resources
type FunctionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []Function `json:"items"`
}
