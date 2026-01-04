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

	v1athenzdomain "github.com/AthenZ/k8s-athenz-syncer/pkg/apis/athenz/v1"

	"github.com/mlajkim/aegis/internal/config"
	"github.com/mlajkim/aegis/internal/syncer"
	"github.com/mlajkim/aegis/pkg/athenz"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

// AthenzDomainController reconciles a AthenzDomain object
type AthenzDomainController struct {
	client.Client
	Scheme       *runtime.Scheme
	Cfg          *config.Config
	SyncerClient *syncer.Syncer
}

// +kubebuilder:rbac:groups=athenz.io,resources=athenzdomains,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;list;watch

func (r *AthenzDomainController) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)

	athenzDomain := &v1athenzdomain.AthenzDomain{}
	if err := r.Get(ctx, req.NamespacedName, athenzDomain); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Full name won't be available; instead, it will only store the parent domain
	fullDomain := athenz.CombineDomains(r.Cfg.Syncer.ParentDomain, athenzDomain.Name)

	// TODO: Make it No Two For Loops
	for _, wantRole := range r.Cfg.Syncer.Roles {
		for _, gotRole := range athenzDomain.Spec.SignedDomain.Domain.Roles {
			if wantRole.Name != athenz.FullRoleNameIntoRoleName(string(gotRole.Name)) {
				continue // i.e) gotRole.Name := eks.users.ajktown-api:role.k8s_ns_viewers
			}

			members := []string{}
			for _, member := range gotRole.Members {
				members = append(members, string(member))
			}

			// Create RoleBindings into k8s namespace:
			r.SyncerClient.ServicesIntoK8sRb(ctx, athenz.DomainIntoNs(r.Cfg.Syncer.ParentDomain, fullDomain), wantRole.Name, members)
		}
	}

	log.Info("Successfully reconciled AthenzDomain", "name", fullDomain)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AthenzDomainController) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1athenzdomain.AthenzDomain{}).
		Complete(r)
}
