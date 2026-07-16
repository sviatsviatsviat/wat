package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapPostToolUse(e sdkclaude.PostToolUse, ev *model.Event) {
	ev.Tool = adapter.NewToolCall(e.ToolName, e.ToolInput.Raw(), e.ToolUseID)
	ev.Result = &model.ToolResult{Raw: adapter.CloneRaw(e.ToolResponse), Text: adapter.RawToText(e.ToolResponse)}
}

func mapPostToolOutput(res model.Result) any {
	var updated any
	if res.UpdatedOutput != nil {
		updated = *res.UpdatedOutput
	}
	if res.Context == "" && updated == nil {
		return nil
	}
	return sdkclaude.BuildPostToolUseOutput(res.Context, updated)
}
