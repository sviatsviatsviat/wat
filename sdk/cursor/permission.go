package cursor

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

// PermissionDecision is a permission verdict label on permission-gating events.
type PermissionDecision string

const (
	// DecisionAllow permits the action.
	DecisionAllow PermissionDecision = "allow"
	// DecisionDeny blocks the action.
	DecisionDeny PermissionDecision = "deny"
	// DecisionAsk escalates to the user.
	DecisionAsk PermissionDecision = "ask"
)

// PermissionOutput is the response for permission-gating events.
// Construct via PermissionResults builders and With* methods. A nil value is a no-op.
type PermissionOutput interface {
	run.Output
	isPermissionOutput()
	// WithUserMessage sets a user-facing message.
	WithUserMessage(msg string) PermissionOutput
	// WithAgentMessage sets an agent-facing message.
	WithAgentMessage(msg string) PermissionOutput
	// WithUpdatedInput replaces tool arguments on preToolUse when set.
	WithUpdatedInput(input map[string]any) PermissionOutput
}

type permissionOutput struct {
	decision     PermissionDecision
	userMessage  string
	agentMessage string
	updatedInput map[string]any
}

func (permissionOutput) isPermissionOutput() {}

// IsZero reports whether this hook response is empty.
func (o permissionOutput) IsZero() bool {
	return o.decision == "" && o.userMessage == "" && o.agentMessage == "" && o.updatedInput == nil
}

// WithUserMessage sets a user-facing message.
func (o permissionOutput) WithUserMessage(msg string) PermissionOutput {
	o.userMessage = msg
	return o
}

// WithAgentMessage sets an agent-facing message.
func (o permissionOutput) WithAgentMessage(msg string) PermissionOutput {
	o.agentMessage = msg
	return o
}

// WithUpdatedInput replaces tool arguments on preToolUse when set.
func (o permissionOutput) WithUpdatedInput(input map[string]any) PermissionOutput {
	o.updatedInput = input
	return o
}

// PermissionResults is the hook-scoped response builder supplied to permission On* handlers by registration.
type PermissionResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionOutput
	// Deny returns a deny verdict with an agent-facing message.
	Deny(agentMessage string) PermissionOutput
	// Ask returns an ask verdict with an agent-facing message.
	Ask(agentMessage string) PermissionOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() PermissionOutput
	isPermissionResults()
}

type permissionResults struct{}

func (permissionResults) isPermissionResults() {}

// Allow returns an allow verdict.
func (permissionResults) Allow() PermissionOutput {
	return permissionOutput{decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing message.
func (permissionResults) Deny(agentMessage string) PermissionOutput {
	return permissionOutput{decision: DecisionDeny, agentMessage: agentMessage}
}

// Ask returns an ask verdict with an agent-facing message.
func (permissionResults) Ask(agentMessage string) PermissionOutput {
	return permissionOutput{decision: DecisionAsk, agentMessage: agentMessage}
}

// Noop returns an empty response (silent stdout).
func (permissionResults) Noop() PermissionOutput {
	return permissionOutput{}
}

type permissionGateResults struct{}

func (permissionGateResults) allow() PermissionOutput {
	return permissionOutput{decision: DecisionAllow}
}

func (permissionGateResults) deny(agentMessage string) PermissionOutput {
	return permissionOutput{decision: DecisionDeny, agentMessage: agentMessage}
}

func (permissionGateResults) ask(agentMessage string) PermissionOutput {
	return permissionOutput{decision: DecisionAsk, agentMessage: agentMessage}
}

func (permissionGateResults) noop() PermissionOutput {
	return permissionOutput{}
}

// Encode renders this output as Cursor stdout JSON.
func (o permissionOutput) Encode() ([]byte, int, error) {
	out := map[string]any{}
	if o.decision != "" {
		out["permission"] = string(o.decision)
	}
	if o.userMessage != "" {
		out["user_message"] = o.userMessage
	}
	if o.agentMessage != "" {
		out["agent_message"] = o.agentMessage
	}
	if o.updatedInput != nil {
		out["updated_input"] = o.updatedInput
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, 0, err
	}
	exitCode := 0
	if o.decision == DecisionDeny {
		exitCode = PermissionDenyExit
	}
	return b, exitCode, nil
}

// Merge combines other into this permission output.
func (o permissionOutput) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(permissionOutput)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	var warnings []string
	decision, agentMessage := hookkit.MergeRankedString(
		string(o.decision), o.agentMessage,
		string(b.decision), b.agentMessage,
		hookkit.PermissionRankString,
	)
	userMessage, w := hookkit.TakeLastString("userMessage", o.userMessage, b.userMessage)
	if w != "" {
		warnings = append(warnings, w)
	}
	updatedInput, w := hookkit.TakeLastMap("updatedInput", o.updatedInput, b.updatedInput)
	if w != "" {
		warnings = append(warnings, w)
	}
	return permissionOutput{
		decision:     PermissionDecision(decision),
		userMessage:  userMessage,
		agentMessage: agentMessage,
		updatedInput: updatedInput,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o permissionOutput) Stop() bool {
	return o.decision == DecisionDeny
}
