package syncer

import (
	"fmt"

	"github.com/mlajkim/aegis/pkg/athenz"
)

// NsIntoAthenzDomain's role is to create:
// - Athenz Domain
// - Necessary roles under the domain above
// for given Kubernetes namespace.
// The Parent domain will be decided by the configuration given as "yaml:syncer.parentDomain
// TODO: Use Solution Template in the future, if possible for quicker/efficient way of adding
func (s *Syncer) NsIntoAthenzDomain(ns string) error {
	// 1. CREATE SUB DOMAIN:
	newDomain := fmt.Sprintf("%s.%s", s.c.Syncer.ParentDomain, athenz.NsIntoDomain(ns))
	if _, err := s.athenzClient.PostSubDomainRecursive(newDomain); err != nil {
		return fmt.Errorf("create subdomain failed: %w, wantDomain %s", err, newDomain)
	}

	// 2. CREATE NECESSARY ROLES
	for _, role := range s.c.Syncer.Roles {
		if err := s.athenzClient.PostRole(newDomain, role.Name, s.c.Syncer.ARoleMembers.IncludeGroup); err != nil {
			return fmt.Errorf("create role %s failed: %w", role.Name, err)
		}
	}

	return nil
}
