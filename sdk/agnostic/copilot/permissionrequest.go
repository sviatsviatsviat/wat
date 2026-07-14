package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapPermissionRequest(e sdkcopilot.PermissionRequest, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input(), "")
}
