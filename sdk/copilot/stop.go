package copilot

import (
	"encoding/json"
)

// StopOutput is the response for agentStop and subagentStop events.
type StopOutput struct {
	// Reason is fed back to the agent as the next instruction.
	Reason string
}

func (o StopOutput) isZero() bool {
	return o.Reason == ""
}

// StopResults is the hook-scoped response builder supplied to Chain handlers by registration.
type StopResults interface {
	// FollowUp blocks completion and feeds reason back to the agent.
	FollowUp(reason string) StopOutput
	isStopResults()
}

type stopResults struct{}

func (stopResults) isStopResults() {}

// FollowUp blocks completion and feeds reason back to the agent.
func (stopResults) FollowUp(reason string) StopOutput {
	return StopOutput{Reason: reason}
}

func encodeStop(o StopOutput) ([]byte, int, error) {
	if o.Reason == "" {
		return nil, 0, nil
	}
	out := map[string]any{
		"decision": "block",
		"reason":   o.Reason,
	}
	b, err := json.Marshal(out)
	return b, 0, err
}
