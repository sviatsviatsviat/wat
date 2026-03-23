package cursorcore

import (
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/helpers"
)

type commonFieldExtractor func(HookDataCommon) string

var commonPlaceholderExtractors = map[string]commonFieldExtractor{
	"CONVERSATION_ID": func(hookData HookDataCommon) string { return hookData.ConversationID },
	"GENERATION_ID":   func(hookData HookDataCommon) string { return hookData.GenerationID },
	"MODEL":           func(hookData HookDataCommon) string { return hookData.Model },
	"HOOK_EVENT_NAME": func(hookData HookDataCommon) string { return hookData.HookEventName },
	"CURSOR_VERSION":  func(hookData HookDataCommon) string { return hookData.CursorVersion },
	"WORKSPACE_ROOTS": func(hookData HookDataCommon) string {
		return helpers.JoinSemicolonSeparatedStrings(hookData.WorkspaceRoots)
	},
	"USER_EMAIL":      func(hookData HookDataCommon) string { return helpers.StringFromPtr(hookData.UserEmail) },
	"TRANSCRIPT_PATH": func(hookData HookDataCommon) string { return helpers.StringFromPtr(hookData.TranscriptPath) },
}

// templateBindingsCommon implements [core.TemplateBindings] for HookDataCommon.
type templateBindingsCommon struct {
	hookData HookDataCommon
}

func newTemplateBindingsCommon(hookData HookDataCommon) core.TemplateBindings {
	return templateBindingsCommon{hookData: hookData}
}

// TemplateValue returns the bound string for placeholderKey and whether it is a known key.
func (bindings templateBindingsCommon) TemplateValue(placeholderKey string) (string, bool) {
	extractString, found := commonPlaceholderExtractors[placeholderKey]
	if !found {
		return "", false
	}
	return extractString(bindings.hookData), true
}
