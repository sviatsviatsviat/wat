package cursor

import (
	"strconv"

	"github.com/sviatsviatsviat/wat/internal/cursor/core"
)

var afterShellExecutionPlaceholderExtractors = map[string]cursorcore.EventFieldExtractor[hookDataAfterShellExecutionFields]{
	"COMMAND":  func(hookData hookDataAfterShellExecutionFields) string { return hookData.Command },
	"OUTPUT":   func(hookData hookDataAfterShellExecutionFields) string { return hookData.Output },
	"DURATION": func(hookData hookDataAfterShellExecutionFields) string {
		return strconv.FormatFloat(float64(hookData.Duration), 'f', -1, 32)
	},
	"SANDBOX": func(hookData hookDataAfterShellExecutionFields) string {
		return strconv.FormatBool(hookData.Sandbox)
	},
}
