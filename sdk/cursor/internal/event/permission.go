package event

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
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
	hookkit.Output
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
	// denyExitZero, when true with DecisionDeny, encodes permission JSON with
	// process exit 0. Used for events such as subagentStart where Cursor applies
	// the JSON permission field and treats exit 2 as a raw-stdout deny message.
	denyExitZero bool
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
	// Deny returns a deny verdict with an agent-facing message and typically
	// PermissionDenyExit. Event-specific builders may use user_message and exit 0.
	Deny(agentMessage string) PermissionOutput
	// Ask returns an ask verdict with an agent-facing message. Enforcement is
	// event-specific: beforeShellExecution and beforeMCPExecution escalate to the
	// user; preToolUse accepts ask but does not enforce it today.
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

// Deny returns a deny verdict with an agent-facing message and PermissionDenyExit.
func (permissionResults) Deny(agentMessage string) PermissionOutput {
	return permissionOutput{decision: DecisionDeny, agentMessage: agentMessage}
}

// Ask returns an ask verdict with an agent-facing message. Cursor enforces ask on
// beforeShellExecution and beforeMCPExecution; see those events' godoc for
// contrast with preToolUse and subagentStart.
func (permissionResults) Ask(agentMessage string) PermissionOutput {
	return permissionOutput{decision: DecisionAsk, agentMessage: agentMessage}
}

// Noop returns an empty response (silent stdout).
func (permissionResults) Noop() PermissionOutput {
	return permissionOutput{}
}

// GateResults is embedded by per-event permission Results builders.
type GateResults struct{}

// Allow returns an allow verdict.
func (GateResults) Allow() PermissionOutput {
	return permissionOutput{decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing message.
func (GateResults) Deny(agentMessage string) PermissionOutput {
	return permissionOutput{decision: DecisionDeny, agentMessage: agentMessage}
}

// Ask returns an ask verdict with an agent-facing message. Cursor enforces ask on
// beforeShellExecution and beforeMCPExecution; see those events' godoc for
// contrast with preToolUse and subagentStart.
func (GateResults) Ask(agentMessage string) PermissionOutput {
	return permissionOutput{decision: DecisionAsk, agentMessage: agentMessage}
}

// Noop returns an empty response (silent stdout).
func (GateResults) Noop() PermissionOutput {
	return permissionOutput{}
}

// DenyUserMessage returns a deny verdict with a user-facing message and process
// exit 0. Cursor's subagentStart schema applies permission from JSON and does not
// use agent_message; exit 2 would re-wrap stdout as the user message.
func (GateResults) DenyUserMessage(userMessage string) PermissionOutput {
	return permissionOutput{
		decision:     DecisionDeny,
		userMessage:  userMessage,
		denyExitZero: true,
	}
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
	if o.decision == DecisionDeny && !o.denyExitZero {
		exitCode = PermissionDenyExit
	}
	return b, exitCode, nil
}

// Merge combines other into this permission output.
func (o permissionOutput) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
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
		denyExitZero: o.denyExitZero || b.denyExitZero,
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o permissionOutput) Stop() bool {
	return o.decision == DecisionDeny
}

// NewPermissionResults returns a PermissionResults builder for handlers outside this package.
func NewPermissionResults() PermissionResults {
	return permissionResults{}
}
