package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapSessionStart(e sdkcursor.SessionStart, ev *model.Event) {
	ev.Life = &model.Lifecycle{Model: e.Model, Background: e.IsBackgroundAgent}
}

func mapSessionStartOutput(res model.Result) any {
	if res.Context == "" {
		return nil
	}
	return sdkcursor.BuildSessionStartOutput(res.Context)
}
