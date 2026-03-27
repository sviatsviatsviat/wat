package run

import (
	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
	"github.com/sviatsviatsviat/wat/internal/helpers"
	"github.com/sviatsviatsviat/wat/internal/template"
)

type commonFieldExtractor func(cursorcore.HookDataCommon) string

var commonPlaceholderExtractors = map[string]commonFieldExtractor{
	"CONVERSATION_ID": func(hookData cursorcore.HookDataCommon) string { return hookData.ConversationID },
	"GENERATION_ID":   func(hookData cursorcore.HookDataCommon) string { return hookData.GenerationID },
	"MODEL":           func(hookData cursorcore.HookDataCommon) string { return hookData.Model },
	"HOOK_EVENT_NAME": func(hookData cursorcore.HookDataCommon) string { return hookData.HookEventName },
	"CURSOR_VERSION":  func(hookData cursorcore.HookDataCommon) string { return hookData.CursorVersion },
	"WORKSPACE_ROOTS": func(hookData cursorcore.HookDataCommon) string {
		return helpers.JoinSemicolonSeparatedStrings(hookData.WorkspaceRoots)
	},
	"USER_EMAIL":      func(hookData cursorcore.HookDataCommon) string { return helpers.StringFromPtr(hookData.UserEmail) },
	"TRANSCRIPT_PATH": func(hookData cursorcore.HookDataCommon) string { return helpers.StringFromPtr(hookData.TranscriptPath) },
}

type templateBindingsCommon struct {
	hookData cursorcore.HookDataCommon
}

func newTemplateBindingsCommon(hookData cursorcore.HookDataCommon) template.TemplateBindings {
	return templateBindingsCommon{hookData: hookData}
}

func (bindings templateBindingsCommon) TemplateValue(placeholderKey string) (string, bool) {
	extractString, found := commonPlaceholderExtractors[placeholderKey]
	if !found {
		return "", false
	}
	return extractString(bindings.hookData), true
}
