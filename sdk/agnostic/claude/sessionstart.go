package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapSessionStart maps a Claude SessionStart hook into a unified Event.
func MapSessionStart(e sdkclaude.SessionStart, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSessionStart)
	ev.Life = &model.Lifecycle{Source: e.Source, Model: e.Model}
	return ev
}
