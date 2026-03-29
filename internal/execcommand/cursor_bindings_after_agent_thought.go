package execcommand

import (
	"strconv"

	"github.com/sviatsviatsviat/wat/internal/cursor"
)

var afterAgentThoughtPlaceholderExtractors = map[string]eventFieldExtractor[cursor.AfterAgentThoughtFields]{
	"DURATION_MS": func(hookData cursor.AfterAgentThoughtFields) string {
		return strconv.FormatInt(hookData.DurationMs, 10)
	},
}
