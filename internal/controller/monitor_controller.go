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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	logr "sigs.k8s.io/controller-runtime/pkg/log"

	uptimekumav1alpha1 "github.com/alemuro/uptime-kuma-k8s/api/v1alpha1"
	"github.com/alemuro/uptime-kuma-k8s/internal/uptimekumaapi"
)

// MonitorReconciler reconciles a Monitor object
type MonitorReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	UptimeAPI uptimekumaapi.UptimeKumaAPI
}

//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=monitors,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=monitors/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=uptimekuma.aleix.cloud,resources=monitors/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Monitor object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *MonitorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	log := logr.Log.WithValues("monitor", req.NamespacedName)

	// Check if the monitor already exists, if not create a new one
	monitor := &uptimekumav1alpha1.Monitor{}
	err := r.Get(ctx, req.NamespacedName, monitor)
	if err != nil && apierrors.IsNotFound(err) {
		// If the custom resource is not found then, it usually means that it was deleted or not created
		// In this way, we will stop the reconciliation

		// TODO: Delete logic
		log.Info("monitor %s/%s not found\n", req.Namespace, req.Name)
		r.UptimeAPI.DeleteMonitor(req.Name)
		return ctrl.Result{}, nil
	} else if err != nil {
		fmt.Println(err, "Error while getting monitor")
		return ctrl.Result{}, err
	}

	log.Info("Reconciling monitor", req.Namespace, req.Name)
	r.UptimeAPI.CreateMonitor(monitor.Name, monitor.Spec.URL, monitor.Spec.Interval, monitor.Spec.Tags)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&uptimekumav1alpha1.Monitor{}).
		Complete(r)
}
