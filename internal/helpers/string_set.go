package helpers

import "sort"

// StringSetSortedKeys returns the keys of set sorted ascending; empty set returns nil.
func StringSetSortedKeys(set map[string]struct{}) []string {
	if len(set) == 0 {
		return nil
	}

	out := make([]string, 0, len(set))
	for key := range set {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}
