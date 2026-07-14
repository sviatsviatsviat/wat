package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStart is the subagentStart hook event.
type SubagentStart struct {
	Envelope
	// AgentName is the agent name (VS Code).
	AgentName string `json:"agent_name"`
	// AgentNameCamel is the agent name (camelCase).
	AgentNameCamel string `json:"agentName"`
	// AgentDisplayName is the display name (VS Code).
	AgentDisplayName string `json:"agent_display_name"`
	// AgentDisplayNameCamel is the display name (camelCase).
	AgentDisplayNameCamel string `json:"agentDisplayName"`
	// AgentDescription is the agent description (camelCase).
	AgentDescription string `json:"agentDescription"`
}

// EventName returns the canonical hook event name.
func (SubagentStart) EventName() string { return EventSubagentStart }

// Name returns the agent name from either wire format.
func (e SubagentStart) Name() string {
	if e.AgentName != "" {
		return e.AgentName
	}
	return e.AgentNameCamel
}

// DisplayName returns the agent display name from either wire format.
func (e SubagentStart) DisplayName() string {
	if e.AgentDisplayName != "" {
		return e.AgentDisplayName
	}
	return e.AgentDisplayNameCamel
}

// SubagentStartOutput is the response for subagentStart events.
type SubagentStartOutput struct {
	// AdditionalContext injects model context.
	AdditionalContext string
}

func (o SubagentStartOutput) isZero() bool {
	return o.AdditionalContext == ""
}

// SubagentStartResults is the hook-scoped response builder supplied to Chain handlers by registration.
type SubagentStartResults interface {
	// Context returns a context-injection-only SubagentStart result.
	Context(text string) SubagentStartOutput
	isSubagentStartResults()
}

type subagentStartResults struct{}

func (subagentStartResults) isSubagentStartResults() {}

// Context returns a context-injection-only SubagentStart result.
func (subagentStartResults) Context(text string) SubagentStartOutput {
	return SubagentStartOutput{AdditionalContext: text}
}

func init() {
	registerDecoder(EventSubagentStart, decodeAs[SubagentStart])
}

// SubagentStart registers a SubagentStart handler.
func (c *Chain) SubagentStart(fn func(context.Context, SubagentStartHook, SubagentStartResults) (SubagentStartOutput, error)) *Chain {
	if fn == nil {
		return c
	}
	registerHandler(func(ctx context.Context, ev SubagentStart) (SubagentStartOutput, error) {
		return fn(ctx, NewHook(run.InvocationFrom(ctx), ev), subagentStartResults{})
	})
	return &Chain{}
}
