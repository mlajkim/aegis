package athenz

import (
	"fmt"
	"io"
)

// curl -k -X PUT "https://localhost:4443/zms/v1/domain/eks.users.ajktown-api/role/k8s_ns_admins" \
// 	--cert ./athenz_distribution/certs/athenz_admin.cert.pem \
// 	--key ./athenz_distribution/keys/athenz_admin.private.pem \
// 	-H "Content-Type: application/json" \
// 	-d '{
// 		"name": "k8s_ns_admins"
// 	}'

// TODO: Maybe my knowledge for POSt and PUT is not correct and Athenz is using it correctly
// TODO: that case I should modify the function name to PutRole, to emphasize it will work EVEN when role exists
// PostRole creates a new role under the specified parent domain.
// Of course, this is not meant for TLD, where creating TLD is only for Athenz administrators.
func (c *AthenzClient) PostRole(domain, newRole string, expand bool) error {
	// if role exists, we should not do something here, because the creation sets members to empty:
	if _, err := c.GetRole(domain, newRole, expand); err == nil {
		return nil
	}

	// They use PUT for role creation, but we are using the PostRole name for consistency:
	// Also it returns 204!
	resp, err := c.Put("/domain/"+domain+"/role/"+newRole, map[string]interface{}{
		"name": newRole,
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("athenz error - status: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
