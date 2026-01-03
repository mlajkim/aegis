package config

import "time"

type Config struct {
	Athenz Athenz `yaml:"athenz"`
	Syncer Syncer `yaml:"syncer"`
}

type Athenz struct {
	ZmsURL           string `yaml:"zmsUrl"`
	CertPath         string `yaml:"certPath"`
	KeyPath          string `yaml:"keyPath"`
	DomainCrdVersion string `yaml:"domainCrdVersion"`
}

type Syncer struct {
	// Shared:
	ParentDomain string       `yaml:"parentDomain"`
	UserTld      string       `yaml:"userTld"`
	Roles        []RoleConfig `yaml:"roles"`

	// Specific:
	ExcludedNamespacesRaw []string            `yaml:"excludedNamespaces"` // Raw
	ExcludedNamespaces    map[string]struct{} `yaml:"-"`                  // Processed
	ARoleMembers          ARoleMembers        `yaml:"athenzRoleMembers"`
	Namespace             Namespaces          `yaml:"namespace"`
}

type RoleConfig struct {
	Name  string       `yaml:"name"`  // Name of the role both in kubernetes & Athenz
	Rules []PolicyRule `yaml:"rules"` // K8s RBAC Policy Rules applied to the role named above
}

type PolicyRule struct {
	APIGroups []string `yaml:"apiGroups"`
	Resources []string `yaml:"resources"`
	Verbs     []string `yaml:"verbs"`
}

type ARoleMembers struct {
	Interval     time.Duration `yaml:"interval"`
	IncludeGroup bool          `yaml:"includeGroup"`
}

type Namespaces struct {
}
