package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

// MapNotification maps a Copilot Notification hook into a unified Event.
func MapNotification(e sdkcopilot.Notification, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindNotification)
	ev.Note = &model.Note{Type: e.NotificationType, Title: e.Title, Message: e.Message}
	return ev
}
