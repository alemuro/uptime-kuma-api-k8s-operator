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

	uptimekumav1alpha1 "github.com/alemuro/uptime-kuma-k8s/api/v1alpha1"
	"github.com/alemuro/uptime-kuma-k8s/internal/uptimekumaapi"

	networkingv1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logr "sigs.k8s.io/controller-runtime/pkg/log"
)

// IngressReconciler reconciles a Namespace object
type IngressReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	UptimeAPI uptimekumaapi.UptimeKumaAPI
}

//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=tags/finalizers,verbs=update
//+kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logr.FromContext(ctx)

	ingress := &networkingv1.Ingress{}
	err := r.Get(ctx, req.NamespacedName, ingress)

	monitor := uptimekumav1alpha1.Monitor{
		ObjectMeta: ctrl.ObjectMeta{
			Name:      req.Name,
			Namespace: "default",
		},
		Spec: uptimekumav1alpha1.MonitorSpec{
			URL:      fmt.Sprintf("https://%s", ingress.Spec.Rules[0].Host),
			Interval: 60,
			Tags:     []string{fmt.Sprintf("k8s-%s", req.Name)},
		},
	}

	// Remove monitor
	if err != nil && apierrors.IsNotFound(err) {
		fmt.Println("Ingress not found, deleting monitor", req.Name)
		r.Delete(ctx, &monitor)
		return ctrl.Result{}, err
	} else if err != nil {
		fmt.Println(err, "Error while getting ingress")
		return ctrl.Result{}, err
	}

	// Create monitor
	err = r.Create(ctx, &monitor)
	if err != nil && apierrors.IsAlreadyExists(err) {
		fmt.Println("Montior already exists", req.Name)
	} else if err != nil {
		fmt.Println(err, "Error while creating monitor")
	} else {
		fmt.Println("Reconciling ingress", req.Name)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		Complete(r)
}
