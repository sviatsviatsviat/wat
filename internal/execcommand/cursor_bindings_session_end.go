package execcommand

import (
	"strconv"

	"github.com/sviatsviatsviat/wat/internal/cursor"
)

var sessionEndPlaceholderExtractors = map[string]eventFieldExtractor[cursor.SessionEndFields]{
	"SESSION_ID": func(hookData cursor.SessionEndFields) string {
		return hookData.SessionID
	},
	"REASON": func(hookData cursor.SessionEndFields) string {
		return hookData.Reason
	},
	"DURATION_MS": func(hookData cursor.SessionEndFields) string {
		return strconv.FormatInt(hookData.DurationMs, 10)
	},
	"IS_BACKGROUND": func(hookData cursor.SessionEndFields) string {
		return strconv.FormatBool(hookData.IsBackgroundAgent)
	},
	"FINAL_STATUS": func(hookData cursor.SessionEndFields) string {
		return hookData.FinalStatus
	},
	"ERROR_MESSAGE": func(hookData cursor.SessionEndFields) string {
		return hookData.ErrorMessage
	},
}
