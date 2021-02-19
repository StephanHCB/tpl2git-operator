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

	// the source

	// repo url for the blueprint
	BlueprintRepoUrl string `json:"blueprint_repo_url,omitempty"`
	// branch to read the blueprint from, defaults to the main branch if left unset
	BlueprintBranch string `json:"blueprint_branch,omitempty"`
	// specifies to read 'generator-<blueprint_name>.yaml', selects which blueprint to generate, defaults to 'main' if left unset
	BlueprintName string `json:"blueprint_name,omitempty"`

	// the target

	// repo url to render to if the spec has changed
	TargetRepoUrl string `json:"target_repo_url,omitempty"`
	// branch to commit to if the spec has changed (will be created if does not exist or updated), defaults to 'update' if left unset
	TargetBranch string `json:"target_branch,omitempty"`
	// branch to fork from if the target branch does not yet exist, defaults to the main branch if left unset
	TargetBranchForkFrom string `json:"target_branch_fork_from,omitempty"`
	// filename of the spec file that is placed in the render output, defaults to 'generated-main.yaml' if left unset
	TargetSpecFile string `json:"target_spec_file,omitempty"`

	// the actual parameter values to be set when performing the render operation. These are written into the target_spec_file
	//
	// changing these is what triggers the render operation because the resource becomes out of sync
	Parameters map[string]string `json:"parameters"`
}

// RendererStatus defines the observed state of Renderer
type RendererStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// the parameter values used during last render
	//
	// this allows the operator to determine if there is anything to be done
	// (and gives us some debugging info)
	CurrentParameters map[string]string `json:"current_parameters"`
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
