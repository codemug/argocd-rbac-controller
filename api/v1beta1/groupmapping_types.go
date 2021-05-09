/*
Copyright 2021.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GroupMappingSpec defines the desired state of GroupMapping
type GroupMappingSpec struct {
	Mappings []MappingSpec `json:"mappings,omitempty,omitempty"`
}

type MappingSpec struct {
	// The name of the group to be created
	GroupName string `json:"groupName,omitempty"`
	// The name of the role to map this group on
	RoleName string `json:"roleName,omitempty"`
}

// GroupMappingStatus defines the observed state of GroupMapping
type GroupMappingStatus struct {
	Status  Status `json:"status,omitempty"`
	Details string `json:"details,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster

// GroupMapping is the Schema for the groupmappings API
type GroupMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GroupMappingSpec   `json:"spec,omitempty"`
	Status GroupMappingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GroupMappingList contains a list of GroupMapping
type GroupMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GroupMapping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GroupMapping{}, &GroupMappingList{})
}
