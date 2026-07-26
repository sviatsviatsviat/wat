package precompact

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
)

// Event is the preCompact hook event.
// Cursor documents this as observational: handlers cannot block or modify
// compaction, but may return a user_message via Results.UserMessage.
type Event struct {
	event.Envelope
	// Trigger is the compaction trigger ("auto" or "manual").
	Trigger string `json:"trigger"`
	// ContextUsagePercent is the current context window usage (0-100).
	ContextUsagePercent int `json:"context_usage_percent"`
	// ContextTokens is the current context window token count.
	ContextTokens int `json:"context_tokens"`
	// ContextWindowSize is the maximum context window size in tokens.
	ContextWindowSize int `json:"context_window_size"`
	// MessageCount is the number of messages in the conversation.
	MessageCount int `json:"message_count"`
	// MessagesToCompact is the number of messages that will be summarized.
	MessagesToCompact int `json:"messages_to_compact"`
	// IsFirstCompaction reports whether this is the first compaction for the conversation.
	IsFirstCompaction bool `json:"is_first_compaction"`
}

// EventName returns the canonical hook event name.
func (Event) EventName() string { return event.PreCompact }

// register registers this hook event decoder on c.
func register(c *hookkit.Codec) {
	c.Register(event.PreCompact, hookkit.EventDecoder[Event](c))
}
