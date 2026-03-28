package execcommand

import (
	"github.com/sviatsviatsviat/wat/internal/cursor"
	"github.com/sviatsviatsviat/wat/internal/helpers"
)

type commonFieldExtractor func(cursor.HookDataCommon) string

var commonPlaceholderExtractors = map[string]commonFieldExtractor{
	"CONVERSATION_ID": func(hookData cursor.HookDataCommon) string { return hookData.ConversationID },
	"GENERATION_ID":   func(hookData cursor.HookDataCommon) string { return hookData.GenerationID },
	"MODEL":           func(hookData cursor.HookDataCommon) string { return hookData.Model },
	"HOOK_EVENT_NAME": func(hookData cursor.HookDataCommon) string { return hookData.HookEventName },
	"CURSOR_VERSION":  func(hookData cursor.HookDataCommon) string { return hookData.CursorVersion },
	"WORKSPACE_ROOTS": func(hookData cursor.HookDataCommon) string {
		return helpers.JoinSemicolonSeparatedStrings(hookData.WorkspaceRoots)
	},
	"USER_EMAIL":      func(hookData cursor.HookDataCommon) string { return helpers.StringFromPtr(hookData.UserEmail) },
	"TRANSCRIPT_PATH": func(hookData cursor.HookDataCommon) string { return helpers.StringFromPtr(hookData.TranscriptPath) },
}

type templateBindingsCommon struct {
	hookData cursor.HookDataCommon
}

func newTemplateBindingsCommon(hookData cursor.HookDataCommon) templateBindings {
	return templateBindingsCommon{hookData: hookData}
}

func (bindings templateBindingsCommon) TemplateValue(placeholderKey string) (string, bool) {
	extractString, found := commonPlaceholderExtractors[placeholderKey]
	if !found {
		return "", false
	}
	return extractString(bindings.hookData), true
}
