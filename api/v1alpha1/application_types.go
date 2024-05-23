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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// +kubebuilder:object:generate=true
type AppSpec struct {
	Image         string            `json:"image"`
	ContainerPort int32             `json:"containerPort"`
	Replicas      *int32            `json:"replicas,omitempty"`
	AppType       string            `json:"appType,omitempty"` // back, front-spa, front-srr
	Annotations   map[string]string `json:"annotations,omitempty"`
	ContainerName string            `json:"containerName"`
	// +optional
	LifeCycle     *v1.Lifecycle		`json:"lifeCycle,omitempty"`
}

type PodDisruptionBudgetSpec struct {
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
	// +optional
	MinAvailable   *int32 `json:"minAvailable,omitempty"`
	MaxUnavailable *int32 `json:"maxUnavailable,omitempty"`
}

type ServiceAccountSpec struct {
	Create *bool `json:"create,omitempty"`
	// +optional
	Annotations map[string]string `json:"annotations,omitempty"`
	// +optional
	AutomountServiceAccountToken *bool `json:"automountServiceAccountToken,omitempty"`
}

type ServiceSpec struct {
	Enabled     *bool             `json:"enabled,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
}

type IngressSpec struct {
	Enabled     bool              `json:"enabled"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Rules       IngressSpecRules  `json:"rules"`
}

// +kubebuilder:object:generate=true
type IngressSpecRules struct {
	Host string `json:"host,omitempty"`
	// +listType=atomic
	Paths []IngressPath `json:"paths"`
}

type IngressPath struct {
	// +optional
	Path string `json:"path,omitempty"`
	// +optional
	ServiceName string `json:"serviceName,omitempty"`
	// +optional
	Port *int32 `json:"port,omitempty"`
}

type ProbeSpec struct {
	Startup   *v1.Probe `json:"startup,omitempty"`
	Liveness  *v1.Probe `json:"liveness,omitempty"`
	Readiness *v1.Probe `json:"readiness,omitempty"`
}

type SchedulerSpec struct {
	NodeSelector            map[string]string       `json:"nodeSelector,omitempty"`
	PodDisruptionBudgetSpec PodDisruptionBudgetSpec `json:"podDisruptionBudget,omitempty"`
	// TODO: need hpa spec
	// TODO: need affinity spec
}

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	App       AppSpec       `json:"app"`
	Scheduler SchedulerSpec `json:"scheduler,omitempty"`
	Probe                         ProbeSpec   `json:"probe,omitempty"`
	TerminationGracePeriodSeconds *int64      `json:"terminationGracePeriodSeconds,omitempty"`
	Service                       ServiceSpec `json:"service,omitempty"`
	Ingress                       IngressSpec `json:"ingress,omitempty"`
}

// ApplicationStatus defines the observed state of Application
type ApplicationStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ApplicationSpec   `json:"spec,omitempty"`
	Status            ApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
