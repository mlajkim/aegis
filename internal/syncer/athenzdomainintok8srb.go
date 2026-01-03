package syncer

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
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

		for _, wantRole := range s.c.Syncer.Roles {
			// get Athenz Role Members from Athenz:
			users, err := s.athenzClient.GetRoleUserMembers(subDomain, wantRole.Name, s.c.Syncer.ARoleMembers.IncludeGroup)
			if err != nil {
				// Continuing is quite important so that users do not experience their permissions suddenly vanishing
				continue // For maximum resilience, we continue even on errors
			}

			if err := s.ServicesIntoK8sRb(ctx, ns, wantRole.Name, users); err != nil {
				return fmt.Errorf("failed to sync Athenz role '%s' members into k8s RoleBinding in ns '%s': %w", wantRole.Name, ns, err)
			}
		}
	}

	return nil
}
