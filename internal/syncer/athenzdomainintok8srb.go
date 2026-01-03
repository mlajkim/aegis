package syncer

import (
	"context"
	"fmt"
	"reflect"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// AthenzDomainIntoK8sRb (RB = RoleBinding)
func (s *Syncer) AthenzDomainIntoK8sRb(ctx context.Context) error {
	// i.e) subDomains=["eks.users.ajktown-api", "gke.users.ajktown-fe", ...]
	subDomains, err := s.athenzClient.GetSubDomains(s.c.Syncer.ParentDomain)
	if err != nil {
		return err
	}

	// Get current namespaces
	nsList := &corev1.NamespaceList{}
	if err := s.k.List(ctx, nsList); err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	// Build a map of existing namespaces for quick lookup:
	existingNamespaces := make(map[string]bool)
	for _, ns := range nsList.Items {
		existingNamespaces[ns.Name] = true
	}

	for _, subDomain := range subDomains {
		ns := s.athenzClient.GetLeaf(subDomain)
		// !WARNING!
		// ! This operator is not designed to manage any k8s namespaces defined in excludedNamespaces,
		// ! And therefore should not do anything, EVEN when athenz server returns a certain namespaces
		// ! like "kube-system", for example. because if we allowed so,
		// ! it would try to add RoleBindings into "kube-system" namespace,
		// ! AND users inside Athenz roles would get permissions in "kube-system" namespace,
		// ! which is definitely NOT what we want, so we make sure to skip them with "continue":
		if _, excludedNs := s.c.Syncer.ExcludedNamespaces[ns]; excludedNs {
			continue
		}

		// if no such ns found, simply skip it. The job of this syncer is NOT to create namespaces,
		// and expect other controllers to create namespaces as needed:
		if !existingNamespaces[ns] {
			continue
		}

		for _, wantRole := range s.c.Syncer.Roles {
			// get Athenz Role Members from Athenz:
			users, err := s.athenzClient.GetRoleUserMembers(subDomain, wantRole.AthenzRole, s.c.Syncer.ARoleMembers.IncludeGroup)
			if err != nil {
				continue // For maximum resilience, we continue even on errors
			}

			// build subject names in rbacv1 native type from users:
			var subjects []rbacv1.Subject
			for _, user := range users {
				subjects = append(subjects, rbacv1.Subject{
					Kind:     "User",
					Name:     user,
					APIGroup: "rbac.authorization.k8s.io",
				})
			}

			// define the desired RoleBinding object:
			wantRb := &rbacv1.RoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name:      s.buildRoleBindingName(ns, wantRole.AthenzRole),
					Namespace: ns,
					Labels:    map[string]string{"managed-by": "athenz-syncer"},
				},
				Subjects: subjects,
				RoleRef: rbacv1.RoleRef{
					APIGroup: "rbac.authorization.k8s.io",
					Kind:     "Role",
					Name:     s.buildRoleName(ns, wantRole.AthenzRole),
				},
			}

			// Check current rb and create if not exists:
			gotRb := &rbacv1.RoleBinding{}
			if err := s.k.Get(ctx, client.ObjectKeyFromObject(wantRb), gotRb); err != nil {
				if errors.IsNotFound(err) {
					if err := s.k.Create(ctx, wantRb); err != nil {
						return err // not expected to fail creating RoleBinding in Kubernetes
					}
					continue // must continue so that other rolebindings are created!
				}
				return err // not expected error
			}

			// Check if existing subjects are different from desired subjects:
			if !reflect.DeepEqual(gotRb.Subjects, wantRb.Subjects) {
				gotRb.Subjects = wantRb.Subjects // Apply desired subjects
				if err := s.k.Update(ctx, gotRb); err != nil {
					return err // not expected to fail updating RoleBinding in Kubernetes
				}
				continue // proceed to next role after update
			}

			// Keep looping until the end of wantRoles...
		}
	}

	return nil
}
