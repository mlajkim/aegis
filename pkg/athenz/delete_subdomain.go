package athenz

import (
	"fmt"
)

// PostSubDomain creates a new subdomain under the specified top-level domain (tld).
// Of course, this is not meant for TLD, where creating TLD is only for Athenz administrators.
// Please note that you cannot delete TLD using this function. (Athenz will return error)
func (c *AthenzClient) DeleteDomain(domain string) error {
	_, parent, leaf := SplitDomain(domain)

	if _, err := c.Delete("/subdomain/" + parent + "/" + leaf); err != nil {
		return fmt.Errorf("failed to delete subdomain %s: %w", domain, err)
	}

	return nil
}
