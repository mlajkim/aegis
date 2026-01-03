package athenz

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type AthenzMember struct {
	MemberName       string `json:"memberName"`
	Approved         bool   `json:"approved"`
	AuditRef         string `json:"auditRef"`
	RequestPrincipal string `json:"requestPrincipal"`
}

type GetRoleResponse struct {
	Name        string `json:"name"`
	Modified    string `json:"modified"`
	RoleMembers []struct {
		AthenzMember
	} `json:"roleMembers"`
}

// TIP: If there is no one as member, it will NOT return roleMembers field at all.
// curl -k -X GET "https://localhost:4443/zms/v1/domain/eks.users.ajktown-api/role/k8s_ns_admins?expand=true" \
//   --cert ./athenz_distribution/certs/athenz_admin.cert.pem \
//   --key ./athenz_distribution/keys/athenz_admin.private.pem
// {"name":"eks.users.ajktown-api:role.k8s_ns_admins","modified":"2025-12-28T05:31:16.815Z","roleMembers":[{"memberName":"user.mlajkim","approved":true,"auditRef":"added using Athenz UI","requestPrincipal":"user.athenz_admin"}]}

// https://github.com/AthenZ/athenz/blob/master/core/zms/src/main/rdl/Role.rdli#L38-L53

// GetRole returns role information in a given domain and role (modified date)
func (c *AthenzClient) GetRole(domainName, roleName string, expand bool) (GetRoleResponse, error) {
	endpoint := fmt.Sprintf("domain/%s/role/%s", domainName, roleName)

	params := url.Values{}
	if expand {
		params.Add("expand", "true")
	}

	resp, err := c.Get(endpoint, params)
	if err != nil {
		return GetRoleResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return GetRoleResponse{}, fmt.Errorf("athenz error - status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	// JSON Parsing
	var res GetRoleResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return GetRoleResponse{}, fmt.Errorf("failed to parse json: %w", err)
	}

	return res, nil
}
