/*
Copyright 2026.

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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// NamespaceReconciler reconciles a Namespace object
type AthenzDomainController struct {
	client.Client
	Scheme       *runtime.Scheme
	Cfg          *config.Config
	SyncerClient *syncer.Syncer
}

// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=namespaces/finalizers,verbs=update

func (r *AthenzDomainController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	// TODO: Add Logic Here:
	// Isee the AthenzDomain CRD HERE! Log it out:
	athenzDomain := &unstructured.Unstructured{}
	athenzDomain.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "athenz.io", // TODO: We need some kind of static SSOT for these values
		Version: r.Cfg.Athenz.DomainCrdVersion,
		Kind:    "AthenzDomain", // TODO: We need some kind of static SSOT for these values
	})

	if err := r.Get(ctx, req.NamespacedName, athenzDomain); err != nil {
		log.Error(err, "failed to get AthenzDomain CRD instance")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	} else {
		log.Info("Successfully found AthenzDomain CRD", "AthenzDomain CRD Name", req.NamespacedName)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AthenzDomainController) SetupWithManager(mgr ctrl.Manager) error {
	// Define what this AthenzDomainController watches:
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "athenz.io", // TODO: We need some kind of static SSOT for these values
		Version: r.Cfg.Athenz.DomainCrdVersion,
		Kind:    "AthenzDomain", // TODO: We need some kind of static SSOT for these values
	})

	return ctrl.NewControllerManagedBy(mgr).
		For(u).     // Watch AthenzDomain CRD
		Complete(r) // Build the controller
}
