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

package controller

import (
	"context"
	"fmt"

	"github.com/alemuro/uptime-kuma-k8s/internal/uptimekumaapi"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logr "sigs.k8s.io/controller-runtime/pkg/log"

	uptimekumav1alpha1 "github.com/alemuro/uptime-kuma-k8s/api/v1alpha1"
)

// TagReconciler reconciles a Tag object
type TagReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	UptimeAPI uptimekumaapi.UptimeKumaAPI
}

//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *TagReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logr.FromContext(ctx)

	log := logr.Log.WithValues("tag", req.NamespacedName)

	// Check if the tag already exists, if not create a new one
	tag := &uptimekumav1alpha1.Tag{}
	err := r.Get(ctx, req.NamespacedName, tag)
	if err != nil && apierrors.IsNotFound(err) {
		log.Info("Tag %s/%s not found\n", req.Namespace, req.Name)
		r.UptimeAPI.DeleteTag(req.Name)
		return ctrl.Result{}, nil
	} else if err != nil {
		fmt.Println(err, "Error while getting tag")
		return ctrl.Result{}, err
	}

	log.Info("Reconciling Tag", req.Namespace, req.Name)
	r.UptimeAPI.CreateTag(tag.Name, tag.Spec.Color)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TagReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uptimekumav1alpha1.Tag{}).
		Complete(r)
}
