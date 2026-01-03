package athenz

// GetLeaf simply returns the leaf role for the given full domain name.
// i.e_ eks.users.ajktown-api -> returns leaf domain name "ajktown-api"
// i.e) eks => returns leaf domain name "eks" (Although eks is TLD, it blindly returns the leaf)
func (c *AthenzClient) GetLeaf(domain string) string {
	lastDotIdx := -1
	for i := len(domain) - 1; i >= 0; i-- {
		if domain[i] == '.' {
			lastDotIdx = i
			break
		}
	}

	if lastDotIdx == -1 {
		return domain
	}

	return domain[lastDotIdx+1:]
}
