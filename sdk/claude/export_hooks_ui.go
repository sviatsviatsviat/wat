package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui/messagedisplay"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui/notification"
)

// Notification is the Notification hook event.
type Notification = notification.Event

// NotificationResults is the hook-scoped response builder for Notification.
type NotificationResults = notification.Results

// MessageDisplay is the MessageDisplay hook event.
type MessageDisplay = messagedisplay.Event

// MessageDisplayOutput is the response for MessageDisplay events.
type MessageDisplayOutput = messagedisplay.Output

// MessageDisplayResults is the hook-scoped response builder for MessageDisplay.
type MessageDisplayResults = messagedisplay.Results
