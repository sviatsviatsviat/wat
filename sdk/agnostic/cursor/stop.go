package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapStop(e sdkcursor.Stop, ev *model.Event) {
	ev.Turn = &model.TurnEnd{Status: e.Status, LoopCount: e.LoopCount}
}
