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

package controllers

import (
	"context"
	"github.com/codemug/argocd-rbac-controller/core"
	"k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	argocdv1beta1 "github.com/codemug/argocd-rbac-controller/api/v1beta1"
)

// RoleMappingReconciler reconciles a RoleMapping object
type RoleMappingReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	RbacManager *core.RbacManager
}

// +kubebuilder:rbac:groups=argocd.codemug.io,resources=rolemappings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=argocd.codemug.io,resources=rolemappings/status,verbs=get;update;patch

func (r *RoleMappingReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("argocdgroupmapping", req.NamespacedName)
	log.Info("Reconciling ArgoCdRoleMapping")

	instance := &argocdv1beta1.RoleMapping{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("ArgoCdRoleMapping removed, clearing its rules")
			r.RbacManager.ClearRoleMapping(&ctrl.ObjectMeta{Name: req.Name, Namespace: req.Namespace})
			err := r.RbacManager.Commit(false)
			if err != nil {
				log.Error(err, "Could not save changes")
				return reconcile.Result{}, err
			}
			log.Info("Reconcile complete")
			return reconcile.Result{}, nil
		}
		log.Error(err, "Could not fetch ArgoCdRoleMapping resource")
		return reconcile.Result{}, err
	}

	log.Info("Resource created or updated, applying its rules")
	err = r.RbacManager.ApplyRoleMapping(instance)
	if err != nil {
		log.Error(err, "Could not apply ArgoCdRoleMapping")
		return ctrl.Result{}, err
	}
	if r.RbacManager.IsDirty() {
		err = r.RbacManager.Commit(false)
		if err != nil {
			log.Error(err, "Could not save changes")
			return reconcile.Result{}, err
		}
		err = r.setStatus(ctx, instance, SUCCESS, "Group rules applied")
		if err != nil {
			log.Error(err, "Could not set status")
			return reconcile.Result{}, err
		}
	}

	log.Info("Reconcile complete")
	return ctrl.Result{}, nil
}

func (r *RoleMappingReconciler) setStatus(ctx context.Context, mapping *argocdv1beta1.RoleMapping, status argocdv1beta1.Status, details string) error {
	if mapping.Status.Status != status || mapping.Status.Details != details {
		mapping.Status.Status = status
		mapping.Status.Details = details
		return r.Client.Status().Update(ctx, mapping)
	}
	return nil
}

func (r *RoleMappingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argocdv1beta1.RoleMapping{}).
		Complete(r)
}
