package syncer

import (
	"github.com/mlajkim/aegis/internal/config"
	"github.com/mlajkim/aegis/pkg/athenz"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Syncer struct {
	c            *config.Config
	k            client.Client
	athenzClient *athenz.AthenzClient
}

func New(cfg *config.Config, k client.Client, athenzClient *athenz.AthenzClient) *Syncer {
	return &Syncer{
		c:            cfg,
		k:            k,
		athenzClient: athenzClient,
	}
}

// i.e) if roleName is "dev-role" and parent domain is "example.domain",
// the returned value is "example.domain:role.dev-role"
func (s *Syncer) buildRoleName(ns, roleName string) string {
	return s.c.Syncer.ParentDomain + "." + ns + ":role." + roleName
}

// i.e) if roleName is "dev-role" and parent domain is "example.domain",
// the returned value is "example.domain:role.dev-role:members"
func (s *Syncer) buildRoleBindingName(ns, roleName string) string {
	return s.buildRoleName(ns, roleName) + ":members"
}
