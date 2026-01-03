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

// AthenzSubjectIntoK8sRb syncs RoleBinding for a specific Athenz subject (subDomain)
// ServicesIntoK8sRb syncs a raw list of services (members) into a specific K8s RoleBinding
func (s *Syncer) ServicesIntoK8sRb(ctx context.Context, ns, role string, services []string) error {
	// !WARNING!
	// ! This operator is not designed to manage any k8s namespaces defined in excludedNamespaces,
	// ! And therefore should not do anything, EVEN when athenz server returns a certain namespaces
	// ! like "kube-system", for example. because if we allowed so,
	// ! it would try to add RoleBindings into "kube-system" namespace,
	// ! AND users inside Athenz roles would get permissions in "kube-system" namespace,
	// ! which is definitely NOT what we want, so we make sure to skip them with "continue":
	if _, excludedNs := s.c.Syncer.ExcludedNamespaces[ns]; excludedNs {
		return nil
	}

	// 2. Check if the namespace exists
	// We check existence first to avoid trying to create RoleBindings in non-existent namespaces
	namespace := &corev1.Namespace{}
	if err := s.k.Get(ctx, client.ObjectKey{Name: ns}, namespace); err != nil {
		if errors.IsNotFound(err) {
			return nil // Namespace doesn't exist, nothing to do; it can just wait until namespace is created
		}
		return fmt.Errorf("failed to get namespace %s: %w", ns, err)
	}

	// 3. Build rbacv1.Subject list from the services slice
	var subjects []rbacv1.Subject
	for _, svc := range services {
		subjects = append(subjects, rbacv1.Subject{
			Kind:     "User", // Assuming Athenz services/members map to K8s Users
			Name:     svc,
			APIGroup: "rbac.authorization.k8s.io",
		})
	}

	// 4. Define the desired RoleBinding object
	wantRb := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.buildRoleBindingName(ns, role),
			Namespace: ns,
			Labels:    map[string]string{"managed-by": "athenz-syncer"},
		},
		Subjects: subjects,
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Role", // Or "ClusterRole" depending on your policy
			Name:     s.buildRoleName(ns, role),
		},
	}

	// 5. Reconcile (Create or Update)
	gotRb := &rbacv1.RoleBinding{}
	if err := s.k.Get(ctx, client.ObjectKeyFromObject(wantRb), gotRb); err != nil {
		if errors.IsNotFound(err) {
			// Create if missing
			if err := s.k.Create(ctx, wantRb); err != nil {
				return fmt.Errorf("failed to create RoleBinding %s in %s: %w", wantRb.Name, ns, err)
			}
			return nil
		}
		return err // Unexpected error
	}

	// 6. Update if subjects differ
	if !reflect.DeepEqual(gotRb.Subjects, wantRb.Subjects) {
		gotRb.Subjects = wantRb.Subjects
		if err := s.k.Update(ctx, gotRb); err != nil {
			return fmt.Errorf("failed to update RoleBinding %s in %s: %w", wantRb.Name, ns, err)
		}
	}

	return nil
}
