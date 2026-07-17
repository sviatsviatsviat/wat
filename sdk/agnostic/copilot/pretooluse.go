package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapPreToolUse maps a Copilot PreToolUse hook into a unified Event.
func MapPreToolUse(e sdkcopilot.PreToolUse, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindPreTool)
	ev.Tool = adapter.NewToolCall(e.NativeToolName(), e.Input().Raw(), "")
	return ev
}
