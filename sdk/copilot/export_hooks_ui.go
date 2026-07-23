package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/ui/notification"
)

// Notification is the notification hook event.
type Notification = notification.Event

// NotificationOutput is the response for notification events.
type NotificationOutput = notification.Output

// NotificationResults is the hook-scoped response builder for Notification.
type NotificationResults = notification.Results
