package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Notification is the notification hook event.
type Notification struct {
	Envelope
	// Message is the notification message.
	Message string `json:"message"`
	// Title is the notification title.
	Title string `json:"title"`
	// NotificationType is the notification category (VS Code).
	NotificationType string `json:"notification_type"`
}

// EventName returns the canonical hook event name.
func (Notification) EventName() string { return EventNotification }

// NotificationOutput is the response for notification events.
type NotificationOutput struct {
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o NotificationOutput) isZero() bool {
	return o.AdditionalContext == ""
}

// NotificationResults is the hook-scoped response builder supplied to Chain handlers by registration.
type NotificationResults interface {
	// Context returns a context-injection-only Notification result.
	Context(text string) NotificationOutput
	isNotificationResults()
}

type notificationResults struct{}

func (notificationResults) isNotificationResults() {}

// Context returns a context-injection-only Notification result.
func (notificationResults) Context(text string) NotificationOutput {
	return NotificationOutput{AdditionalContext: text}
}

func init() {
	registerDecoder(EventNotification, decodeAs[Notification])
}

// Notification registers a Notification handler.
func (c *Chain) Notification(fn func(context.Context, NotificationHook, NotificationResults) (NotificationOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev Notification) (NotificationOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), notificationResults{})
	})
	return &Chain{}
}
