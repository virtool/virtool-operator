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
	"k8s.io/apimachinery/pkg/util/intstr"
)

// VirtoolAppSpec defines the desired state of VirtoolApp
type VirtoolAppSpec struct {
	// Components is a map of component names to their desired specifications
	Components map[string]ComponentSpec `json:"components"`

	// UpdateStrategy defines how updates should be performed
	UpdateStrategy UpdateStrategy `json:"updateStrategy,omitempty"`

	// GlobalConfig defines global settings for all components
	GlobalConfig GlobalConfig `json:"globalConfig,omitempty"`
}

// ComponentSpec defines the specification for a single component
type ComponentSpec struct {
	// Version is the desired version of the component
	Version string `json:"version"`

	// Image is the container image for the component
	Image string `json:"image"`

	// Replicas is the desired number of replicas for the component
	Replicas *int32 `json:"replicas,omitempty"`

	// Resources defines the resource requirements for the component
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`

	// UpdateOrder defines the order in which this component should be updated
	UpdateOrder int `json:"updateOrder,omitempty"`

	// PreUpdateJob defines a job to run before updating this component
	PreUpdateJob *Job `json:"preUpdateJob,omitempty"`

	// PostUpdateJob defines a job to run after updating this component
	PostUpdateJob *Job `json:"postUpdateJob,omitempty"`
}

// Job defines a job to be run as part of the update process
type Job struct {
	Image   string   `json:"image"`
	Command []string `json:"command,omitempty"`
	Args    []string `json:"args,omitempty"`
}

// UpdateStrategy defines how updates should be performed
type UpdateStrategy struct {
	// Type is the type of update strategy (e.g., "RollingUpdate", "Recreate")
	Type string `json:"type"`

	// MaxUnavailable defines the maximum number of pods that can be unavailable during the update
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`

	// MaxSurge defines the maximum number of pods that can be created over the desired number of pods
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty"`
}

// GlobalConfig defines global settings for all components
type GlobalConfig struct {
	// Registry is the default container registry to use
	Registry string `json:"registry,omitempty"`

	// ImagePullSecret is the name of the secret containing registry credentials
	ImagePullSecret string `json:"imagePullSecret,omitempty"`

	// Tolerations defines the tolerations for all components
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`

	// NodeSelector defines the node selector for all components
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
}

// VirtoolAppStatus defines the observed state of VirtoolApp
type VirtoolAppStatus struct {
	// ComponentStatus tracks the status of individual components
	ComponentStatus map[string]ComponentStatus `json:"componentStatus"`

	// Conditions represent the latest available observations of the VirtoolApp's state
	Conditions []metav1.Condition `json:"conditions"`

	// LastUpdateTime is the last time the status was updated
	LastUpdateTime metav1.Time `json:"lastUpdateTime"`
}

// ComponentStatus tracks the status of an individual component
type ComponentStatus struct {
	// CurrentVersion is the current version of the component
	CurrentVersion string `json:"currentVersion"`

	// DesiredVersion is the version the component is being updated to
	DesiredVersion string `json:"desiredVersion"`

	// Status is the current status of the component
	Status string `json:"status"`

	// Message provides additional information about the component's status
	Message string `json:"message,omitempty"`

	// ReadyReplicas is the number of replicas that are ready
	ReadyReplicas int32 `json:"readyReplicas"`

	// UpdatedReplicas is the number of replicas that have been updated
	UpdatedReplicas int32 `json:"updatedReplicas"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// VirtoolApp is the Schema for the virtoolapps API
type VirtoolApp struct {
	Spec              VirtoolAppSpec   `json:"spec,omitempty"`
	Status            VirtoolAppStatus `json:"status,omitempty"`
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
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
