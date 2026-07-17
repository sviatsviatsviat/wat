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
// Construct via NotificationResults builders. A nil value is a no-op.
type NotificationOutput interface {
	isNotificationOutput()
}

type notificationOutput struct {
	additionalContext string
}

func (notificationOutput) isNotificationOutput() {}

func (o notificationOutput) isZero() bool {
	return o.additionalContext == ""
}

// NotificationResults is the hook-scoped response builder supplied to Chain handlers by registration.
type NotificationResults interface {
	// Context returns a context-injection-only Notification result.
	Context(text string) NotificationOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() NotificationOutput
	isNotificationResults()
}

type notificationResults struct{}

func (notificationResults) isNotificationResults() {}

// Context returns a context-injection-only Notification result.
func (notificationResults) Context(text string) NotificationOutput {
	return notificationOutput{additionalContext: text}
}

// Noop returns an empty response (silent stdout).
func (notificationResults) Noop() NotificationOutput {
	return notificationOutput{}
}

func (notificationOutput) allowedEvents() []string {
	return []string{EventNotification}
}

func (o notificationOutput) encode() ([]byte, int, error) {
	return encodeAdditionalContext(o.additionalContext)
}

func init() {
	registerDecoder(EventNotification, decodeAs[Notification])
}

// Notification registers a Notification handler.
func (c *Chain) Notification(fn func(context.Context, Hook[Notification], NotificationResults) (NotificationOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(c.registerOwner(), func(ctx context.Context, ev Notification) (NotificationOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), notificationResults{})
	})
	return c
}
