package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapPermissionRequest maps a Copilot PermissionRequest hook into a unified Event.
func MapPermissionRequest(e sdkcopilot.PermissionRequest, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPermissionRequest)
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input().Raw(), "")
	return ev
}
