package athenz

import "strings"

// SplitDomain splits fullDomain into tld, parent, and leaf.
// For example, "eks.users.ajktown-api" becomes:
//
//   - tld: "eks"
//   - parent: "eks.users"
//   - leaf: "ajktown-api"
func SplitDomain(fullDomain string) (tld string, parent string, leaf string) {
	idx := strings.LastIndex(fullDomain, ".")
	if idx == -1 {
		return fullDomain, "", fullDomain
	}

	leaf = fullDomain[idx+1:]

	parent = fullDomain[:idx]

	if firstIdx := strings.Index(fullDomain, "."); firstIdx != -1 {
		tld = fullDomain[:firstIdx]
	} else {
		tld = fullDomain
	}

	return tld, parent, leaf
}

// NsIntoDomain converts Kubernetes namespace name into Athenz domain name
// i.e) "eks-users-ajktown-api" becomes "eks.users.ajktown.api"
func NsIntoDomain(ns string) string {
	return strings.ReplaceAll(ns, "-", ".")
}

// DomainIntoNs converts Athenz domain name into Kubernetes namespace name
// It requires parentDomain to strip the parent domain part.
func DomainIntoNs(parentDomain, domain string) string {
	leaf := strings.TrimPrefix(domain, parentDomain+".")
	return strings.ReplaceAll(leaf, ".", "-")
}

func FullRoleNameIntoRoleName(fullRoleName string) string {
	parts := strings.SplitN(fullRoleName, ":role.", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return fullRoleName
}

func CombineDomains(domain1, domain2 string) string {
	if domain1 == "" {
		return domain2
	}
	if domain2 == "" {
		return domain1
	}
	return domain1 + "." + domain2
}
