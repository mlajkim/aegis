/*
Copyright 2025.

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

	"github.com/mlajkim/aegis/internal/config"
	"github.com/mlajkim/aegis/internal/syncer"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// NamespaceReconciler reconciles a Namespace object
type NamespaceReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	Cfg          *config.Config
	SyncerClient *syncer.Syncer
}

// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=namespaces/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Namespace object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
func (r *NamespaceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	ns := corev1.Namespace{}
	if err := r.Get(ctx, req.NamespacedName, &ns); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if _, exists := r.Cfg.Syncer.ExcludedNamespaces[req.Name]; exists {
		log.V(1).Info("Namespace is excluded from syncer, skipping", "excludedNamespace", req.Name)
		return ctrl.Result{}, nil
	}

	if !ns.DeletionTimestamp.IsZero() {
		// Just simple log is fine for now:
		log.V(1).Info("Namespace is being deleted, skipping", "namespace", ns.Name)
		return ctrl.Result{}, nil
	}

	if err := r.SyncerClient.NsIntoAthenzDomain(ctx, req.Name); err != nil {
		return ctrl.Result{}, err
	}

	if err := r.SyncerClient.NsIntoK8sRole(ctx, req.Name); err != nil {
		return ctrl.Result{}, err
	}

	// Log success:
	log.Info("Successfully reconciled", "namespace", req.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *NamespaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Named("namespace").

		// If you want some control of concurrency:
		// WithEventFilter(predicate.Funcs{
		// 	CreateFunc:  func(e event.CreateEvent) bool { return true },
		// 	DeleteFunc:  func(e event.DeleteEvent) bool { return false }, // Ignore deletions
		// 	UpdateFunc:  func(e event.UpdateEvent) bool { return false }, // Ignore updates
		// 	GenericFunc: func(e event.GenericEvent) bool { return false },
		// }).
		Complete(r)
}
