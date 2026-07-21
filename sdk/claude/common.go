package claude

import "github.com/sviatsviatsviat/wat/sdk/run"

// PermissionDecision is a pre-tool permission verdict label.
type PermissionDecision string

const (
	// DecisionAllow permits the tool call.
	DecisionAllow PermissionDecision = "allow"
	// DecisionDeny blocks the tool call.
	DecisionDeny PermissionDecision = "deny"
	// DecisionAsk escalates to the user.
	DecisionAsk PermissionDecision = "ask"
	// DecisionDefer defers the permission decision.
	DecisionDefer PermissionDecision = "defer"
)

// common holds output fields shared across Claude Code hook responses.
type common struct {
	cont             *bool
	stopReason       string
	suppressOutput   bool
	systemMessage    string
	terminalSequence string
}

// IsZero reports whether this hook response is empty.
func (c common) IsZero() bool {
	return c.cont == nil && c.stopReason == "" && !c.suppressOutput &&
		c.systemMessage == "" && c.terminalSequence == ""
}

// WithContinue sets whether Claude should continue the session.
// Pass false to stop Claude entirely.
func (c common) WithContinue(v bool) common {
	c.cont = &v
	return c
}

// WithStopReason explains why the session was stopped.
func (c common) WithStopReason(reason string) common {
	c.stopReason = reason
	return c
}

// WithSuppressOutput suppresses hook stdout when true.
func (c common) WithSuppressOutput(v bool) common {
	c.suppressOutput = v
	return c
}

// WithSystemMessage sets a user-visible system message.
func (c common) WithSystemMessage(msg string) common {
	c.systemMessage = msg
	return c
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (c common) WithTerminalSequence(seq string) common {
	c.terminalSequence = seq
	return c
}

// CommonOutput is a shared-fields-only response for events that only accept those fields.
// Construct via Results builders and With* methods. A nil value is a no-op.
type CommonOutput interface {
	run.Output
	isCommonOutput()
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) CommonOutput
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) CommonOutput
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) CommonOutput
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) CommonOutput
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) CommonOutput
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) CommonOutput
}

type commonOutput struct {
	common
	eventName         string
	additionalContext string
}

func (commonOutput) isCommonOutput() {}

// IsZero reports whether this hook response is empty.
func (o commonOutput) IsZero() bool {
	return o.common.IsZero() && o.additionalContext == ""
}

// WithAdditionalContext injects model context.
func (o commonOutput) WithAdditionalContext(text string) CommonOutput {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o commonOutput) WithContinue(v bool) CommonOutput {
	o.common = o.common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o commonOutput) WithStopReason(reason string) CommonOutput {
	o.common = o.common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o commonOutput) WithSuppressOutput(v bool) CommonOutput {
	o.common = o.common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o commonOutput) WithSystemMessage(msg string) CommonOutput {
	o.common = o.common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o commonOutput) WithTerminalSequence(seq string) CommonOutput {
	o.common = o.common.WithTerminalSequence(seq)
	return o
}

func applyCommon(top map[string]any, c common) {
	if c.cont != nil {
		top["continue"] = *c.cont
		if !*c.cont && c.stopReason != "" {
			top["stopReason"] = c.stopReason
		}
	}
	if c.systemMessage != "" {
		top["systemMessage"] = c.systemMessage
	}
	if c.suppressOutput {
		top["suppressOutput"] = true
	}
	if c.terminalSequence != "" {
		top["terminalSequence"] = c.terminalSequence
	}
}

func (o commonOutput) encodeInto(top, hso map[string]any) {
	applyCommon(top, o.common)
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o commonOutput) Encode() ([]byte, int, error) {
	return marshalHookOutput(o.eventName, o.encodeInto)
}
