package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

// MapNotification maps a Claude Notification hook into a unified Event.
func MapNotification(e sdkclaude.Notification, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindNotification)
	ev.Note = &model.Note{Type: e.NotificationType, Message: e.Message}
	return ev
}
