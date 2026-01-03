package util

// StrArrayIntoUniqSet converts a string array into a map with unique keys.
// Please note that the "set" is used, as go does not have native set data structure.
func StrArrayIntoUniqSet(given []string) map[string]struct{} {
	uniqSet := make(map[string]struct{})
	for _, item := range given {
		uniqSet[item] = struct{}{}
	}
	return uniqSet
}
