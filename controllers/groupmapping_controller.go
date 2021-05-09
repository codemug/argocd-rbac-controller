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

const (
	SUCCESS = "Success"
)

// GroupMappingReconciler reconciles a GroupMapping object
type GroupMappingReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	RbacManager *core.RbacManager
}

//+kubebuilder:rbac:groups=argocd.codemug.io,resources=groupmappings,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argocd.codemug.io,resources=groupmappings/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argocd.codemug.io,resources=groupmappings/finalizers,verbs=update
func (r *GroupMappingReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("argocdgroupmapping", req.NamespacedName)
	log.Info("Reconciling ArgoCdGroupMapping")
	instance := &argocdv1beta1.GroupMapping{}
	err := r.Client.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("ArgoCdGroupMapping removed, clearing its rules")
			r.RbacManager.ClearGroupMapping(&ctrl.ObjectMeta{Name: req.Name, Namespace: req.Namespace})
			err := r.RbacManager.Commit(false)
			if err != nil {
				log.Error(err, "Could not save changes after removal")
				return reconcile.Result{}, err
			}
			log.Info("Reconcile complete")
			return reconcile.Result{}, nil
		}
		log.Error(err, "Could not fetch ArgoCdGroupMapping resource")
		return reconcile.Result{}, err
	}

	log.Info("Resource created or updated, applying its rules")
	r.RbacManager.ApplyGroupMapping(instance)
	if r.RbacManager.IsDirty() {
		err = r.RbacManager.Commit(false)
		if err != nil {
			log.Error(err, "Could not save changes after update")
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

func (r *GroupMappingReconciler) setStatus(ctx context.Context, mapping *argocdv1beta1.GroupMapping, status argocdv1beta1.Status, details string) error {
	if mapping.Status.Status != status || mapping.Status.Details != details {
		mapping.Status.Status = status
		mapping.Status.Details = details
		return r.Client.Status().Update(ctx, mapping)
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *GroupMappingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argocdv1beta1.GroupMapping{}).
		Complete(r)
}
