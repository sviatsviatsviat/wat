// Package template expands __KEY__ placeholders in argv tokens using [core.TemplateBindings].
package template

import (
	"regexp"

	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/helpers"
)

var placeholderRE = regexp.MustCompile(`__([A-Z_]+)__`)

// RenderTokens replaces each __KEY__ substring in tokens using bindings.
// It returns the rendered argv slice and placeholder keys that were not defined by bindings, sorted.
func RenderTokens(tokens []string, bindings core.TemplateBindings) ([]string, []string) {
	unknownKeysSet := map[string]struct{}{}
	renderedTokens := make([]string, 0, len(tokens))

	for _, argvToken := range tokens {
		renderedArgvToken := placeholderRE.ReplaceAllStringFunc(argvToken, func(fullPlaceholder string) string {
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
		renderedTokens = append(renderedTokens, renderedArgvToken)
	}

	unknownKeysSorted := helpers.StringSetSortedKeys(unknownKeysSet)
	return renderedTokens, unknownKeysSorted
}
