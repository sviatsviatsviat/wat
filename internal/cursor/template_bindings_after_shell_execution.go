package cursor

import (
	"strconv"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// afterShellExecutionFieldExtractor returns one template field from hookDataAfterShellExecution.
type afterShellExecutionFieldExtractor func(hookDataAfterShellExecutionFields) string

var afterShellExecutionPlaceholderExtractors = map[string]afterShellExecutionFieldExtractor{
	"COMMAND":  func(hookData hookDataAfterShellExecutionFields) string { return hookData.Command },
	"OUTPUT":   func(hookData hookDataAfterShellExecutionFields) string { return hookData.Output },
	"DURATION": func(hookData hookDataAfterShellExecutionFields) string {
		return strconv.FormatFloat(float64(hookData.Duration), 'f', -1, 32)
	},
	"SANDBOX": func(hookData hookDataAfterShellExecutionFields) string {
		if hookData.Sandbox {
			return "true"
		}
		return "false"
	},
}

// templateBindingsAfterShellExecution implements [core.TemplateBindings] for afterShellExecution payloads.
type templateBindingsAfterShellExecution struct {
	commonBindings core.TemplateBindings
	hookData       hookDataAfterShellExecutionFields
}

func newTemplateBindingsAfterShellExecution(hookData hookDataAfterShellExecution) core.TemplateBindings {
	return templateBindingsAfterShellExecution{
		commonBindings: newTemplateBindingsCommon(hookData.hookDataCommon),
		hookData:       hookData.hookDataAfterShellExecutionFields,
	}
}

// TemplateValue returns the bound string for placeholderKey and whether it is a known key.
func (bindings templateBindingsAfterShellExecution) TemplateValue(placeholderKey string) (string, bool) {
	extractString, found := afterShellExecutionPlaceholderExtractors[placeholderKey]
	if found {
		return extractString(bindings.hookData), true
	}
	return bindings.commonBindings.TemplateValue(placeholderKey)
}
