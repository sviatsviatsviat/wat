package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapUserPromptSubmit(e sdkclaude.UserPromptSubmit, ev *model.Event) {
	ev.Prompt = e.Prompt
}
