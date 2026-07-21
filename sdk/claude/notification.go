package claude

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Notification is the Notification hook event.
type Notification struct {
	Envelope
	// Message is the notification message.
	Message string `json:"message"`
	// NotificationType is the notification category.
	NotificationType string `json:"notification_type"`
}

// EventName returns the hook event name.
func (Notification) EventName() string { return EventNotification }

func init() {
	codec.Register(EventNotification, hookkit.EventDecoder[Notification](codec))
}

// NotificationResults is the hook-scoped response builder supplied to On* handlers by registration.
type NotificationResults interface {
	// Context returns a context-injection-only Notification result.
	Context(text string) CommonOutput
	isNotificationResults()
}

type notificationResults struct{}

func (notificationResults) isNotificationResults() {}

// Context returns a context-injection-only Notification result.
func (notificationResults) Context(text string) CommonOutput {
	return commonOutput{eventName: EventNotification, additionalContext: text}
}

// Notification registers a Notification handler on the chain.
func (c *chain) Notification(fn func(context.Context, run.Hook[Notification], NotificationResults) (CommonOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[Notification]) (CommonOutput, error) {
		return fn(ctx, hook, notificationResults{})
	}))
	return c
}
