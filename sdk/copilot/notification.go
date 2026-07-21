package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

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
	run.Output
	isNotificationOutput()
}

type notificationOutput struct {
	additionalContext string
}

func (notificationOutput) isNotificationOutput() {}

// IsZero reports whether this hook response is empty.
func (o notificationOutput) IsZero() bool {
	return o.additionalContext == ""
}

// NotificationResults is the hook-scoped response builder supplied to On* handlers by registration.
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

// Encode renders this output as Copilot stdout JSON.
func (o notificationOutput) Encode() ([]byte, int, error) {
	return encodeAdditionalContext(o.additionalContext)
}

func init() {
	codec.Register(EventNotification, hookkit.EventDecoder[Notification](codec))
}

// Notification registers a Notification handler on the chain.
func (c *chain) Notification(fn func(context.Context, run.Hook[Notification], NotificationResults) (NotificationOutput, error)) *chain {
	if fn == nil {
		return c
	}
	c.reg.RegisterHandler(Dialect, run.Handler(func(ctx context.Context, hook run.Hook[Notification]) (NotificationOutput, error) {
		return fn(ctx, hook, notificationResults{})
	}))
	return c
}
