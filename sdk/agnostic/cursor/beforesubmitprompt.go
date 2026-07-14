package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapBeforeSubmitPrompt(e sdkcursor.BeforeSubmitPrompt, ev *model.Event) {
	ev.Prompt = e.Prompt
}
