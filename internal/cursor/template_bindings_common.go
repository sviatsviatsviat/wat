package cursor

import (
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/helpers"
)

// commonFieldExtractor returns one template field from hookDataCommon.
type commonFieldExtractor func(hookDataCommon) string

// commonPlaceholderExtractors maps placeholder keys to field extractors.
var commonPlaceholderExtractors = map[string]commonFieldExtractor{
	"CONVERSATION_ID": func(hookData hookDataCommon) string { return hookData.ConversationID },
	"GENERATION_ID":   func(hookData hookDataCommon) string { return hookData.GenerationID },
	"MODEL":           func(hookData hookDataCommon) string { return hookData.Model },
	"HOOK_EVENT_NAME": func(hookData hookDataCommon) string { return hookData.HookEventName },
	"CURSOR_VERSION":  func(hookData hookDataCommon) string { return hookData.CursorVersion },
	"WORKSPACE_ROOTS": func(hookData hookDataCommon) string {
		return helpers.JoinSemicolonSeparatedStrings(hookData.WorkspaceRoots)
	},
	"USER_EMAIL":      func(hookData hookDataCommon) string { return helpers.StringFromPtr(hookData.UserEmail) },
	"TRANSCRIPT_PATH": func(hookData hookDataCommon) string { return helpers.StringFromPtr(hookData.TranscriptPath) },
}

// templateBindingsCommon implements [core.TemplateBindings] for hookDataCommon.
type templateBindingsCommon struct {
	hookData hookDataCommon
}

// newTemplateBindingsCommon returns [core.TemplateBindings] backed by hookData.
func newTemplateBindingsCommon(hookData hookDataCommon) core.TemplateBindings {
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
