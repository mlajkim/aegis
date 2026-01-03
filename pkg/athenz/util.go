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
