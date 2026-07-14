package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func mapNotification(e sdkcopilot.Notification, ev *model.Event) {
	ev.Note = &model.Note{Type: e.NotificationType, Title: e.Title, Message: e.Message}
}
