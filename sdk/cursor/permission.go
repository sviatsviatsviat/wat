package cursor

import (
	"encoding/json"
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
type PermissionOutput struct {
	// Decision is the permission verdict (allow, deny, ask).
	Decision PermissionDecision
	// UserMessage is shown to the user.
	UserMessage string
	// AgentMessage is shown to the agent.
	AgentMessage string
	// UpdatedInput replaces tool arguments on preToolUse when set.
	UpdatedInput map[string]any
}

func (o PermissionOutput) isZero() bool {
	return o.Decision == "" && o.UserMessage == "" && o.AgentMessage == "" && o.UpdatedInput == nil
}

// PermissionResults is the hook-scoped response builder supplied to permission Chain handlers by registration.
type PermissionResults interface {
	// Allow returns an allow verdict.
	Allow() PermissionOutput
	// Deny returns a deny verdict with an agent-facing message.
	Deny(agentMessage string) PermissionOutput
	isPermissionResults()
}

type permissionResults struct{}

func (permissionResults) isPermissionResults() {}

// Allow returns an allow verdict.
func (permissionResults) Allow() PermissionOutput {
	return PermissionOutput{Decision: DecisionAllow}
}

// Deny returns a deny verdict with an agent-facing message.
func (permissionResults) Deny(agentMessage string) PermissionOutput {
	return PermissionOutput{Decision: DecisionDeny, AgentMessage: agentMessage}
}

type permissionGateResults struct{}

func (permissionGateResults) allow() PermissionOutput {
	return PermissionOutput{Decision: DecisionAllow}
}

func (permissionGateResults) deny(agentMessage string) PermissionOutput {
	return PermissionOutput{Decision: DecisionDeny, AgentMessage: agentMessage}
}

func encodePermission(eventName string, o PermissionOutput) ([]byte, int, error) {
	out := map[string]any{}
	if o.Decision != "" {
		out["permission"] = string(o.Decision)
	}
	if o.UserMessage != "" {
		out["user_message"] = o.UserMessage
	}
	if o.AgentMessage != "" {
		out["agent_message"] = o.AgentMessage
	}
	if o.UpdatedInput != nil && (eventName == "" || eventName == EventPreToolUse) {
		out["updated_input"] = o.UpdatedInput
	}
	if len(out) == 0 {
		return nil, 0, nil
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, 0, err
	}
	exitCode := 0
	if o.Decision == DecisionDeny {
		exitCode = PermissionDenyExit
	}
	return b, exitCode, nil
}
