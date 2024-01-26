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

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logr "sigs.k8s.io/controller-runtime/pkg/log"

	uptimekumav1alpha1 "github.com/alemuro/uptime-kuma-k8s/api/v1alpha1"
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	UptimeAPI uptimekumaapi.UptimeKumaAPI
}

//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logr.FromContext(ctx)

	namespace := &corev1.Namespace{}
	err := r.Get(ctx, req.NamespacedName, namespace)
	if err != nil && apierrors.IsNotFound(err) {
		fmt.Println("Namespace not found")
		tag := uptimekumav1alpha1.Tag{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      fmt.Sprintf("k8s-%s", req.Name),
				Namespace: "default",
			},
			Spec: uptimekumav1alpha1.TagSpec{
				Color: "black",
			},
		}
		r.Delete(ctx, &tag)
		return ctrl.Result{}, nil
	} else if err != nil {
		fmt.Println(err, "Error while getting namespace")
		return ctrl.Result{}, err
	}

	fmt.Println("Reconciling namespace", req.Name)

	tag := uptimekumav1alpha1.Tag{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      fmt.Sprintf("k8s-%s", req.Name),
			Namespace: "default",
		},
		Spec: uptimekumav1alpha1.TagSpec{
			Color: "black",
		},
	}
	err = r.Create(ctx, &tag)
	if err != nil && apierrors.IsAlreadyExists(err) {
		fmt.Println("Tag already exists")
	} else if err != nil {
		fmt.Println(err, "Error while creating tag")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Complete(r)
}
