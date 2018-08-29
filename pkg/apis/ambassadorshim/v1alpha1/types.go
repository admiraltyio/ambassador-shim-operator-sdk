package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type MappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Mapping `json:"items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Mapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              MappingSpec   `json:"spec"`
	Status            MappingStatus `json:"status,omitempty"`
}

type MappingSpec struct {
	Prefix  string `json:"prefix"`
	Service string `json:"service"`
}
type MappingStatus struct {
	Configured bool `json:"configured"`
	UpToDate   bool `json:"upToDate"`
}
