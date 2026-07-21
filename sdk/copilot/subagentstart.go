package copilot

import (
	"context"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// SubagentStart is the SubagentStart hook event.
type SubagentStart struct {
	Envelope
	// AgentName is the agent name.
	AgentName string `json:"agent_name"`
	// AgentDisplayName is the display name.
	AgentDisplayName string `json:"agent_display_name"`
	// AgentDescription is the agent description.
	AgentDescription string `json:"agent_description"`
}

// EventName returns the canonical hook event name.
func (SubagentStart) EventName() string { return EventSubagentStart }

// Name returns the agent name.
func (e SubagentStart) Name() string {
	return e.AgentName
}

// DisplayName returns the agent display name.
func (e SubagentStart) DisplayName() string {
	return e.AgentDisplayName
}

// SubagentStartOutput is the response for SubagentStart events.
// Construct via SubagentStartResults builders. A nil value is a no-op.
type SubagentStartOutput interface {
	run.Output
	isSubagentStartOutput()
}

type subagentStartOutput struct {
	additionalContext string
}

func (subagentStartOutput) isSubagentStartOutput() {}

// IsZero reports whether this hook response is empty.
func (o subagentStartOutput) IsZero() bool {
	return o.additionalContext == ""
}

// SubagentStartResults is the hook-scoped response builder supplied to On* handlers by registration.
type SubagentStartResults interface {
	// Context returns a context-injection-only SubagentStart result.
	Context(text string) SubagentStartOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() SubagentStartOutput
	isSubagentStartResults()
}

type subagentStartResults struct{}

func (subagentStartResults) isSubagentStartResults() {}

// Context returns a context-injection-only SubagentStart result.
func (subagentStartResults) Context(text string) SubagentStartOutput {
	return subagentStartOutput{additionalContext: text}
}

// Noop returns an empty response (silent stdout).
func (subagentStartResults) Noop() SubagentStartOutput {
	return subagentStartOutput{}
}

// Encode renders this output as Copilot stdout JSON.
func (o subagentStartOutput) Encode() ([]byte, int, error) {
	return encodeAdditionalContext(o.additionalContext)
}

func init() {
	codec.Register(EventSubagentStart, hookkit.EventDecoder[SubagentStart](codec))
}

// SubagentStart registers a SubagentStart handler on the chain.
func (c *chain) SubagentStart(fn func(context.Context, run.Hook[SubagentStart], SubagentStartResults) (SubagentStartOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(c.reg, func(ctx context.Context, ev SubagentStart) (SubagentStartOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), subagentStartResults{})
	})
	return c
}
