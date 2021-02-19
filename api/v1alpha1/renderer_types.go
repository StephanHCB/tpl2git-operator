/*
Copyright 2021 StephanHCB.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// RendererSpec defines the desired state of Renderer
type RendererSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Renderer. Edit Renderer_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// RendererStatus defines the observed state of Renderer
type RendererStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Renderer is the Schema for the renderers API
type Renderer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RendererSpec   `json:"spec,omitempty"`
	Status RendererStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RendererList contains a list of Renderer
type RendererList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Renderer `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Renderer{}, &RendererList{})
}
