package run

import (
	"regexp"

	"github.com/sviatsviatsviat/wat/internal/helpers"
)

var placeholderRE = regexp.MustCompile(`__([A-Z_]+)__`)

// templateBindings maps placeholder keys to values for substituting __KEY__ segments in command argument strings.
// TemplateValue reports whether key is defined; value may be empty when ok is true.
type templateBindings interface {
	TemplateValue(key string) (value string, ok bool)
}

// renderTokens replaces each __KEY__ substring in tokens using bindings.
// It returns the rendered argument slice and placeholder keys that were not defined by bindings, sorted.
func renderTokens(tokens []string, bindings templateBindings) ([]string, []string) {
	unknownKeysSet := map[string]struct{}{}
	renderedTokens := make([]string, 0, len(tokens))

	for _, token := range tokens {
		renderedToken := placeholderRE.ReplaceAllStringFunc(token, func(fullPlaceholder string) string {
			m := placeholderRE.FindStringSubmatch(fullPlaceholder)
			if len(m) < 2 {
				return ""
			}
			placeholderKey := m[1]

			replacement, keyDefined := bindings.TemplateValue(placeholderKey)
			if !keyDefined {
				unknownKeysSet[placeholderKey] = struct{}{}
				return ""
			}
			return replacement
		})
		renderedTokens = append(renderedTokens, renderedToken)
	}

	unknownKeysSorted := helpers.StringSetSortedKeys(unknownKeysSet)
	return renderedTokens, unknownKeysSorted
}
