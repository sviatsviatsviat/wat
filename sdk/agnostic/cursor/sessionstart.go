package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapSessionStart(e sdkcursor.SessionStart, ev *model.Event) {
	ev.Life = &model.Lifecycle{Model: e.Model, Background: e.IsBackgroundAgent}
}
