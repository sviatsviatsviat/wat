// Package helpers provides small string and set utilities.
package helpers

import "strings"

// JoinSemicolonSeparatedStrings joins values with ";". Nil or empty slice yields "".
func JoinSemicolonSeparatedStrings(values []string) string {
	return strings.Join(values, ";")
}

// StringFromPtr returns the string pointed to by p, or "" if p is nil.
func StringFromPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
