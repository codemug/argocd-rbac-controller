package core

import (
	"github.com/codemug/argocd-rbac-controller/api/v1beta1"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestRbacManager_ApplyGroupMapping(t *testing.T) {
	manager := NewRbacManager(nil, "my-cm", "my-namespace", "default-policy")
	groupMapping := v1beta1.GroupMapping{
		ObjectMeta: v1.ObjectMeta{Name: "testGroups", Namespace: "test"},
		Spec: v1beta1.GroupMappingSpec{
			Mappings: []v1beta1.MappingSpec{
				{
					GroupName: "myGroup",
					RoleName:  "myRole",
				},
			},
		},
	}
	manager.ApplyGroupMapping(&groupMapping)
	assert.True(t, manager.dirty)
	assert.Contains(t, manager.rules, "groups/test/testGroups")
}
