package athenz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type DomainResponse struct {
	Domains []string `json:"names"`
}

func (c *AthenzClient) GetSubDomains(parentDomain string) ([]string, error) {
	resp, err := c.Get("/domain", url.Values{
		"prefix": []string{parentDomain},
	})

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("athenz error - status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// Json parsing:
	var result DomainResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse json: %w", err)
	}

	// Remove the first item and return:
	return result.Domains[1:], nil
}
