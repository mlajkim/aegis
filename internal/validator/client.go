package validator

import (
	"github.com/mlajkim/aegis/internal/config"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Validator struct {
	c                    *config.Config
	k                    client.Client
	supportedCRDVersions map[string]struct{}
}

func New(cfg *config.Config, k client.Client) *Validator {
	return &Validator{
		c: cfg,
		k: k,
		supportedCRDVersions: map[string]struct{}{
			"v1": {},
		},
	}
}
