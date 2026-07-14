package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapErrorOccurred(e sdkcopilot.ErrorOccurred, ev *model.Event) {
	if detail, ok := e.Detail(); ok {
		noteType := detail.Name
		if noteType == "" {
			noteType = e.Context()
		}
		ev.Note = &model.Note{
			Type:        noteType,
			Message:     detail.Message,
			Recoverable: e.Recoverable,
		}
	}
}
