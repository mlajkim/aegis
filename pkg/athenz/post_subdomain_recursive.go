package athenz

import (
	"fmt"
	"strings"
)

// PostSubDomainRecursive creates the domain hierarchy from top to bottom.
// If the domain is "kr.sports.baseball", it ensures:
// 1. "kr.sports" exists (if not, creates it)
// 2. "kr.sports.baseball" exists (if not, creates it)
//
// It returns an error if the input is a TLD (Top Level Domain), as TLD creation requires system admin privileges.
func (c *AthenzClient) PostSubDomainRecursive(domain string) (*PostSubDomainResponse, error) {
	if res, err := c.GetDomain(domain); err == nil {
		return &PostSubDomainResponse{
			Description:  res.Description,
			Org:          res.Org,
			Name:         res.Name,
			Modified:     res.Modified,
			ID:           res.ID,
			AuditEnabled: res.AuditEnabled,
		}, nil
	}

	// 1. Validation: TLD Creation is not allowed
	if !strings.Contains(domain, ".") {
		return nil, fmt.Errorf("refusing to create TLD '%s': creation is restricted to Athenz System Admins", domain)
	}

	// 2. Split into parts
	// Ex: "kr.sports.baseball" -> ["eks", "users", "ajktown", "api"]
	parts := strings.Split(domain, ".")

	// Start with the TLD (we assume TLD exists)
	currentDomain := parts[0] // i.e) "eks"

	var lastResponse *PostSubDomainResponse
	var err error

	// 3. Iterate from the first subdomain level
	// i=1: "eks" + "." + "users" -> Create "eks.users"
	// i=2: "eks.users" + "." + "ajktown" -> Create "eks.users.ajktown"
	for i := 1; i < len(parts); i++ {
		currentDomain = currentDomain + "." + parts[i]

		// Call the single creation function
		// Note: PostSubDomain internally checks if the domain exists.
		lastResponse, err = c.PostSubDomain(currentDomain)
		if err != nil {
			return nil, fmt.Errorf("failed to ensure domain segment '%s': %w", currentDomain, err)
		}
	}

	// Return the response of the final leaf domain
	return lastResponse, nil
}
