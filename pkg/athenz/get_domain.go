package athenz

import (
	"encoding/json"
	"fmt"
	"io"
)

type GetDomainResponse struct {
	Description  string `json:"description"`
	Org          string `json:"org"`
	AuditEnabled bool   `json:"auditEnabled"`
	Name         string `json:"name"`
	Modified     string `json:"modified"`
	ID           string `json:"id"`
}

// PostSubDomain creates a new subdomain under the specified top-level domain (tld).
// Of course, this is not meant for TLD, where creating TLD is only for Athenz administrators.
func (c *AthenzClient) GetDomain(domain string) (*PostSubDomainResponse, error) {
	resp, err := c.Get("/domain/"+domain, nil)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("athenz error - status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Json parsing:
	var result PostSubDomainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	return &result, nil
}
