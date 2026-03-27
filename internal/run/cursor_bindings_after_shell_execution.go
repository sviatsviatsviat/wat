package run

import (
	"strconv"

	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
)

var afterShellExecutionPlaceholderExtractors = map[string]eventFieldExtractor[cursorcore.AfterShellExecutionFields]{
	"COMMAND": func(hookData cursorcore.AfterShellExecutionFields) string { return hookData.Command },
	"OUTPUT":  func(hookData cursorcore.AfterShellExecutionFields) string { return hookData.Output },
	"DURATION": func(hookData cursorcore.AfterShellExecutionFields) string {
		return strconv.FormatFloat(float64(hookData.Duration), 'f', -1, 32)
	},
	"SANDBOX": func(hookData cursorcore.AfterShellExecutionFields) string {
		return strconv.FormatBool(hookData.Sandbox)
	},
}
