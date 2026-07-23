package event

// PermissionDecision is a pre-tool permission verdict label.
type PermissionDecision string

const (
	// DecisionAllow permits the tool call.
	DecisionAllow PermissionDecision = "allow"
	// DecisionDeny blocks the tool call.
	DecisionDeny PermissionDecision = "deny"
	// DecisionAsk escalates to the user.
	DecisionAsk PermissionDecision = "ask"
)
