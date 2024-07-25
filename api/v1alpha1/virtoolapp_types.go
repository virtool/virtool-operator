/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// VirtoolAppSpec defines the desired state of the application
type VirtoolAppSpec struct {
	// Version is the desired version of the application
	Version string `json:"version"`

	// Components is a list of components for the application
	Components []ComponentSpec `json:"components"`
}

// ComponentSpec defines the specification for a single component
type ComponentSpec struct {
	// Name is the name of the component
	Name string `json:"name"`

	// Image is the container image for the component
	Image string `json:"image"`

	// Replicas is the desired number of replicas for the component
	Replicas int32 `json:"replicas,omitempty"`

	// Resources defines the resource requirements for the component
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// PreUpdateJob defines a job to run before updating this component
	PreUpdateJob *JobSpec `json:"preUpdateJob,omitempty"`

	// PostUpdateJob defines a job to run after updating this component
	PostUpdateJob *JobSpec `json:"postUpdateJob,omitempty"`
}

// JobSpec defines a job to be run as part of the update process
type JobSpec struct {
	Image   string   `json:"image"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// VirtoolAppStatus defines the observed state of the application
type VirtoolAppStatus struct {
	// CurrentVersion is the current version of the application
	CurrentVersion string `json:"currentVersion"`

	// ComponentsStatus tracks the status of individual components
	ComponentsStatus []ComponentStatus `json:"componentsStatus"`

	// Conditions represent the latest available observations of the VirtoolApp's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// ComponentStatus tracks the status of an individual component
type ComponentStatus struct {
	// Name is the name of the component
	Name string `json:"name"`

	// CurrentVersion is the current version of the component
	CurrentVersion string `json:"currentVersion"`

	// Status is the current status of the component
	Status string `json:"status"`

	// ReadyReplicas is the number of replicas that are ready
	ReadyReplicas int32 `json:"readyReplicas"`

	// UpdatedReplicas is the number of replicas that have been updated
	UpdatedReplicas int32 `json:"updatedReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VirtoolApp is the Schema for the virtoolapps API
type VirtoolApp struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtoolAppSpec   `json:"spec,omitempty"`
	Status VirtoolAppStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// VirtoolAppList contains a list of VirtoolApp
type VirtoolAppList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VirtoolApp `json:"items"`
}

func init() {
	SchemeBuilder.Register(&VirtoolApp{}, &VirtoolAppList{})
}
