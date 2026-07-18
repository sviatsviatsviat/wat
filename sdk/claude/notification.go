package claude

import (
	"context"

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
	registerDecoder(EventNotification, decodeAs[Notification])
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
	return commonOutput{additionalContext: text}
}

// OnNotification registers a Notification handler.
func OnNotification(fn func(context.Context, Hook[Notification], NotificationResults) (CommonOutput, error)) *chain {
	return (&chain{}).Notification(fn)
}

// Notification registers another Notification handler on the chain.
func (c *chain) Notification(fn func(context.Context, Hook[Notification], NotificationResults) (CommonOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Notification) (CommonOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), notificationResults{})
	})
	return c
}
