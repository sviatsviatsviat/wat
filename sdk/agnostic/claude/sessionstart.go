package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapSessionStart(e sdkclaude.SessionStart, ev *model.Event) {
	ev.Life = &model.Lifecycle{Source: e.Source, Model: e.Model}
}

func mapSessionStartOutput(res model.Result) any {
	if res.Context == "" {
		return nil
	}
	return sdkclaude.BuildSessionStartOutput(res.Context)
}
