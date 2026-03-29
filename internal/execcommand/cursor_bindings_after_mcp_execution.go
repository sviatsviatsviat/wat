package execcommand

import (
	"strconv"

	"github.com/sviatsviatsviat/wat/internal/cursor"
)

var afterMCPExecutionPlaceholderExtractors = map[string]eventFieldExtractor[cursor.AfterMCPExecutionFields]{
	"TOOL_NAME": func(hookData cursor.AfterMCPExecutionFields) string {
		return hookData.ToolName
	},
	"DURATION": func(hookData cursor.AfterMCPExecutionFields) string {
		return strconv.FormatFloat(hookData.Duration, 'f', -1, 64)
	},
}
