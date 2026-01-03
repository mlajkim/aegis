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
// i.e) "eks.users.ajktown.api" becomes "eks-users-ajktown-api" (Please note thhat if Parent Domain is defined, it should be stripped first)
func DomainIntoNs(domain string) string {
	return strings.ReplaceAll(domain, ".", "-")
}
