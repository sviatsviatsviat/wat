package execcommand

import (
	"strconv"

	"github.com/sviatsviatsviat/wat/internal/cursor"
)

var afterShellExecutionPlaceholderExtractors = map[string]eventFieldExtractor[cursor.AfterShellExecutionFields]{
	"COMMAND": func(hookData cursor.AfterShellExecutionFields) string { return hookData.Command },
	"OUTPUT":  func(hookData cursor.AfterShellExecutionFields) string { return hookData.Output },
	"DURATION": func(hookData cursor.AfterShellExecutionFields) string {
		return strconv.FormatFloat(float64(hookData.Duration), 'f', -1, 32)
	},
	"SANDBOX": func(hookData cursor.AfterShellExecutionFields) string {
		return strconv.FormatBool(hookData.Sandbox)
	},
}
