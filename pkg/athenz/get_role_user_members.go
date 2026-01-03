package athenz

// GetRoleUserMembers returns list of unique user members in a given domain and role,
// The logic depends on the `GetUser()` function to filter only user.* members.`
func (c *AthenzClient) GetRoleUserMembers(domainName, roleName string, expand bool) ([]string, error) {
	res, err := c.GetRole(domainName, roleName, expand)
	if err != nil {
		return nil, err
	}

	var userMembers []string
	for _, member := range res.RoleMembers {
		if member.MemberName[0:len(c.UserTld)+1] == c.UserTld+"." { //Required to distinguish tld "user" vs "useruser".
			userMembers = append(userMembers, member.MemberName)
		}
	}
	return userMembers, nil
}
