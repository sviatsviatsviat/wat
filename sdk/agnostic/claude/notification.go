package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

func mapNotification(e sdkclaude.Notification, ev *model.Event) {
	ev.Note = &model.Note{Type: e.NotificationType, Message: e.Message}
}
