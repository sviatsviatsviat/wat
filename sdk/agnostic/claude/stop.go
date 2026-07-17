package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapStop maps a Claude Stop hook into a unified Event.
func MapStop(e sdkclaude.Stop, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindStop)
	ev.Turn = &model.TurnEnd{StopHookActive: e.StopHookActive, LastAssistantMessage: e.LastAssistantMessage}
	return ev
}
