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

type Status string

const (
	SUCCESS Status = "Success"
	FAILURE Status = "Failure"
)


// ArgoCdGroupMappingSpec defines the desired state of ArgoCdGroupMapping
type ArgoCdGroupMappingSpec struct {
	Mappings []GroupMappingSpec `json:"mappings,omitempty,omitempty"`
}

type GroupMappingSpec struct {
	// The name of the group to be created
	GroupName string `json:"groupName,omitempty"`
	// The name of the role to map this group on
	RoleName  string `json:"roleName,omitempty"`
}

// ArgoCdGroupMappingStatus defines the observed state of ArgoCdGroupMapping
type ArgoCdGroupMappingStatus struct {
	Status  Status `json:"status,omitempty"`
	Details string `json:"details,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ArgoCdGroupMapping is the Schema for the argocdgroupmappings API
type ArgoCdGroupMapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ArgoCdGroupMappingSpec   `json:"spec,omitempty"`
	Status ArgoCdGroupMappingStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ArgoCdGroupMappingList contains a list of ArgoCdGroupMapping
type ArgoCdGroupMappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ArgoCdGroupMapping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ArgoCdGroupMapping{}, &ArgoCdGroupMappingList{})
}
