package core

import (
	"context"
	"errors"
	"fmt"
	"github.com/codemug/argocd-rbac-controller/api/v1beta1"
	v1 "k8s.io/api/core/v1"
	errors2 "k8s.io/apimachinery/pkg/api/errors"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strings"
)

const (
	GROUP_FORMAT                = "g, %s, role:%s"
	PERMISSION_FORMAT           = "p, role:%s, %s, %s, %s, allow"
	DEFAULT_CONFIGMAP_NAME      = "argocd-rbac-cm"
	DEFAULT_CONFIGMAP_NAMESPACE = "argocd"
	DEFAULT_POLICY_KEY          = "policy.default"
	DEFAULT_RBAC_KEY            = "policy.csv"
)

var RESOURCES = map[string]bool{
	"clusters":     true,
	"projects":     true,
	"applications": true,
	"repositories": true,
	"certificates": true,
	"*":            true,
}

var ACTIONS = map[string]bool{
	"get":      true,
	"create":   true,
	"update":   true,
	"delete":   true,
	"sync":     true,
	"override": true,
	"action":   true,
}

type RbacManager struct {
	rules         map[string]map[string]bool
	dirty         bool
	client        client.Client
	configMap     *v1.ConfigMap
	defaultPolicy string
}

func NewRbacManager(client client.Client, configMapName string, configMapNamespace string, defaultPolicy string) RbacManager {
	if configMapName == "" {
		configMapName = DEFAULT_CONFIGMAP_NAME
	}
	if configMapNamespace == "" {
		configMapNamespace = DEFAULT_CONFIGMAP_NAMESPACE
	}
	cm := v1.ConfigMap{ObjectMeta: v12.ObjectMeta{Name: configMapName, Namespace: configMapNamespace}}
	return RbacManager{client: client, configMap: &cm, defaultPolicy: defaultPolicy, rules: make(map[string]map[string]bool)}
}

func (r *RbacManager) ApplyGroupMapping(mapping *v1beta1.GroupMapping) {
	groupRules := make(map[string]bool)
	for _, v := range mapping.Spec.Mappings {
		rule := fmt.Sprintf(GROUP_FORMAT, v.GroupName, v.RoleName)
		groupRules[rule] = true
	}
	if len(groupRules) > 0 {
		r.rules[GetTypedNamespacedName(&mapping.TypeMeta, &mapping.ObjectMeta)] = groupRules
		r.dirty = true
	}
}

func (r *RbacManager) ClearMapping(typeMeta *v12.TypeMeta, objectMeta *v12.ObjectMeta) {
	namespacedName := GetTypedNamespacedName(typeMeta, objectMeta)
	if _, ok := r.rules[namespacedName]; ok {
		delete(r.rules, namespacedName)
		r.dirty = true
	}
}

func (r *RbacManager) ApplyRoleMapping(mapping *v1beta1.RoleMapping) error {
	allPermissions := make(map[string]bool)
	for _, v := range mapping.Spec.Roles {
		permissions, err := r.getPermissions(&v)
		if err != nil {
			return err
		}
		for k, _ := range permissions {
			allPermissions[k] = true
		}
	}
	if len(allPermissions) > 0 {
		r.rules[GetTypedNamespacedName(&mapping.TypeMeta, &mapping.ObjectMeta)] = allPermissions
		r.dirty = true
	}
	return nil
}

func (r *RbacManager) getPermissions(mapping *v1beta1.RoleSpec) (map[string]bool, error) {
	permissions := make(map[string]bool)
	for _, v := range mapping.Permissions {
		if _, ok := RESOURCES[v.Resource]; !ok {
			return nil, errors.New(fmt.Sprintf("invalid resource specified %s", v.Resource))
		}
		if _, ok := ACTIONS[v.Action]; !ok {
			return nil, errors.New(fmt.Sprintf("invalid action specified %s", v.Action))
		}
		if v.Instance == "" {
			v.Instance = "*"
		}
		permissions[fmt.Sprintf(PERMISSION_FORMAT, mapping.Name, v.Resource, v.Action, v.Instance)] = true
	}
	return permissions, nil
}

func (r *RbacManager) getFullRbac() string {
	allRules := make(map[string]bool)
	for _, v := range r.rules {
		for k, _ := range v {
			allRules[k] = true
		}
	}
	keys := make([]string, len(allRules))
	i := 0
	for k, _ := range allRules {
		keys[i] = k
		i++
	}
	return strings.Join(keys, "\n")
}

func (r *RbacManager) Commit(force bool) error {
	if r.dirty || force {
		configMapContents := make(map[string]string)
		configMapContents[DEFAULT_RBAC_KEY] = r.getFullRbac()
		configMapContents[DEFAULT_POLICY_KEY] = r.defaultPolicy
		found := &v1.ConfigMap{}
		err := r.client.Get(context.TODO(), types.NamespacedName{Name: r.configMap.Name, Namespace: r.configMap.Namespace}, found)
		if err != nil {
			if errors2.IsNotFound(err) {
				r.configMap.Data = configMapContents
				err = r.client.Create(context.TODO(), r.configMap)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			found.Data = configMapContents
			err = r.client.Update(context.TODO(), found)
			if err != nil {
				return err
			}
		}
		r.dirty = false
	}
	return nil
}

func GetTypedNamespacedName(typeMeta *v12.TypeMeta, objectMeta *v12.ObjectMeta) string {
	return fmt.Sprintf("%s/%s/%s", typeMeta.Kind, objectMeta.Namespace, objectMeta.Name)
}
