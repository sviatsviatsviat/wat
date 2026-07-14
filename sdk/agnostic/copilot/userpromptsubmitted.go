package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapUserPromptSubmitted(e sdkcopilot.UserPromptSubmitted, ev *model.Event) {
	ev.Prompt = e.Prompt
}
