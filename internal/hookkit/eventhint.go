package hookkit

// EventHint holds runtime decode configuration shared by Cursor and Copilot.
type EventHint struct {
	Hint string
}

// HintOption configures EventHint.
type HintOption func(*EventHint)

// WithEventHint supplies the native event name when hook_event_name is absent.
func WithEventHint(name string) HintOption {
	return func(c *EventHint) {
		c.Hint = name
	}
}

// ApplyHintOptions applies opts to cfg.
func ApplyHintOptions(cfg *EventHint, opts ...HintOption) {
	for _, opt := range opts {
		opt(cfg)
	}
}

// DefaultEventHint returns an empty EventHint.
func DefaultEventHint() EventHint {
	return EventHint{}
}
