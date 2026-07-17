package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapPermissionRequest maps a Claude PermissionRequest hook into a unified Event.
func MapPermissionRequest(e sdkclaude.PermissionRequest, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPermissionRequest)
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	return ev
}
