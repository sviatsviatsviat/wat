package claude

import (
	"context"
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Elicitation is the Elicitation hook event.
type Elicitation struct {
	Envelope
	// ServerName is the MCP server name.
	ServerName string `json:"server_name"`
	// Message is the elicitation message.
	Message string `json:"message"`
	// Schema is the requested input schema JSON.
	Schema json.RawMessage `json:"requested_schema"`
}

// EventName returns the hook event name.
func (Elicitation) EventName() string { return EventElicitation }

func init() {
	codec.Register(EventElicitation, hookkit.EventDecoder[Elicitation](codec))
}

// ElicitationOutput is the response for Elicitation events.
// Construct via ElicitationResults builders and With* methods.
// A nil value is a no-op.
type ElicitationOutput interface {
	Output
	isElicitationOutput()
	// WithContent sets the elicitation response content.
	WithContent(content map[string]any) ElicitationOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) ElicitationOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) ElicitationOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) ElicitationOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) ElicitationOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) ElicitationOutput
}

type elicitationOutput struct {
	common
	action  string
	content map[string]any
}

func (elicitationOutput) isClaudeOutput() {}

func (elicitationOutput) isElicitationOutput() {}
func (o elicitationOutput) isZero() bool {
	return o.common.isZero() && o.action == "" && o.content == nil
}

// WithContent sets the elicitation response content.
func (o elicitationOutput) WithContent(content map[string]any) ElicitationOutput {
	o.content = content
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o elicitationOutput) WithContinue(v bool) ElicitationOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o elicitationOutput) WithStopReason(reason string) ElicitationOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o elicitationOutput) WithSuppressOutput(v bool) ElicitationOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o elicitationOutput) WithSystemMessage(msg string) ElicitationOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o elicitationOutput) WithTerminalSequence(seq string) ElicitationOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

// ElicitationResults is the hook-scoped response builder supplied to On* handlers by registration.
type ElicitationResults interface {
	// Accept returns an accept action result.
	Accept() ElicitationOutput
	// Decline returns a decline action result.
	Decline() ElicitationOutput
	// Cancel returns a cancel action result.
	Cancel() ElicitationOutput
	isElicitationResults()
}

type elicitationResults struct{}

func (elicitationResults) isElicitationResults() {}

// Accept returns an accept action result.
func (elicitationResults) Accept() ElicitationOutput {
	return elicitationOutput{action: "accept"}
}

// Decline returns a decline action result.
func (elicitationResults) Decline() ElicitationOutput {
	return elicitationOutput{action: "decline"}
}

// Cancel returns a cancel action result.
func (elicitationResults) Cancel() ElicitationOutput {
	return elicitationOutput{action: "cancel"}
}

func (elicitationOutput) allowedEvents() []string {
	return []string{EventElicitation}
}

func (o elicitationOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.action != "" {
		hso["action"] = o.action
	}
	if o.content != nil {
		hso["content"] = o.content
	}
}

// Elicitation registers a Elicitation handler on the chain.
func (c *chain) Elicitation(fn func(context.Context, run.Hook[Elicitation], ElicitationResults) (ElicitationOutput, error)) *chain {
	if fn == nil {
		return c
	}
	registerHandler(c.reg, func(ctx context.Context, ev Elicitation) (ElicitationOutput, error) {
		return fn(ctx, run.NewHook(run.InvocationFrom(ctx), ev), elicitationResults{})
	})
	return c
}
