package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapErrorOccurred maps a Copilot ErrorOccurred hook into a unified Event.
func MapErrorOccurred(e sdkcopilot.ErrorOccurred, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindAgentError)
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
	return ev
}
