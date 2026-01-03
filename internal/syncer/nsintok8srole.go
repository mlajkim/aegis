package syncer

import (
	"context"
	"reflect"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NsIntoK8sRole creates K8s native Role with given namespace
// With the given configuration "yaml:syncer.roles"
// It does NOT create rolebindings for isolated responsibilities.

func (s *Syncer) NsIntoK8sRole(ctx context.Context, ns string) error {
	for _, wantRole := range s.c.Syncer.Roles {
		var rules []rbacv1.PolicyRule
		for _, r := range wantRole.Rules {
			rules = append(rules, rbacv1.PolicyRule{
				APIGroups: r.APIGroups,
				Resources: r.Resources,
				Verbs:     r.Verbs,
			})
		}

		wantRole := &rbacv1.Role{ObjectMeta: metav1.ObjectMeta{Name: s.buildRoleName(ns, wantRole.AthenzRole), Namespace: ns}, Rules: rules}
		gotRole := &rbacv1.Role{}

		err := s.k.Get(ctx, client.ObjectKeyFromObject(wantRole), gotRole)
		if err != nil {
			if errors.IsNotFound(err) {
				if err := s.k.Create(ctx, wantRole); err != nil {
					return err
				}
				continue // must continue so that other roles are created!
			}
			return err // any failures other than NotFound
		}

		// if exists, we want to make sure that the correct permissions are applied:
		if !reflect.DeepEqual(gotRole.Rules, wantRole.Rules) {
			gotRole.Rules = wantRole.Rules
			return s.k.Update(ctx, gotRole)
		}

		return nil // exists, already equal; everyone happy.
	}

	return nil
}
