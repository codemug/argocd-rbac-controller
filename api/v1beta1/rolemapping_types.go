/*


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

// RoleMappingSpec defines the desired state of RoleMapping
type RoleMappingSpec struct {
	// Roles to permissions mapping specification
	Roles []RoleSpec `json:"roles,omitempty"`
}

type PermissionSpec struct {
	// The type of resource on which the permission is to be defined.
	// This can be one of (clusters, projects, applications,
	// repositories, certificates, or * for all)
	Resource string `json:"resource,omitempty"`
	// The action that is being permitted on the specified resource.
	// This can be one of (get, create, update, delete, sync,
	// override, action, or * for all)
	Action string `json:"action,omitempty"`
	// If the permission is to be applied to a specific instance of
	// the resource type, the name of that instance is specified here.
	// If this is not specified, the permission is applied to all
	// instances of the resource type
	// +kubebuilder:validation:Optional
	Instance string `json:"instance,omitempty"`
}

type RoleSpec struct {
	// The name of the role
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Optional
	Permissions []PermissionSpec `json:"permissions,omitempty"`
}

// RoleMappingStatus defines the observed state of RoleMapping
type RoleMappingStatus struct {
	Status  Status `json:"status,omitempty"`
	Details string `json:"details,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// RoleMapping is the Schema for the rolemappings API
type RoleMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RoleMappingSpec   `json:"spec,omitempty"`
	Status RoleMappingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RoleMappingList contains a list of RoleMapping
type RoleMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RoleMapping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RoleMapping{}, &RoleMappingList{})
}
