package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapSessionEnd(e sdkcursor.SessionEnd, ev *model.Event) {
	ev.Life = &model.Lifecycle{Reason: e.Reason, Background: e.IsBackgroundAgent}
}
