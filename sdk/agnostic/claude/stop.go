package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapStop(e sdkclaude.Stop, ev *model.Event) {
	ev.Turn = &model.TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
}

func mapStopOutput(res model.Result) any {
	if res.FollowUp == "" {
		return nil
	}
	return sdkclaude.BuildStopOutput(res.FollowUp)
}
