package copilot

import (
	"encoding/json"
)

// StopOutput is the response for agentStop and subagentStop events.
// Construct via StopResults builders. A nil value is a no-op.
type StopOutput interface {
	Output
	isStopOutput()
}

type stopOutput struct {
	reason string
}

func (stopOutput) isCopilotOutput() {}

func (stopOutput) isStopOutput() {}

func (o stopOutput) isZero() bool {
	return o.reason == ""
}

// StopResults is the hook-scoped response builder supplied to On* handlers by registration.
type StopResults interface {
	// FollowUp blocks completion and feeds reason back to the agent.
	FollowUp(reason string) StopOutput
	// Noop returns an empty response (silent stdout). Prefer nil from handlers when not chaining With*.
	Noop() StopOutput
	isStopResults()
}

type stopResults struct{}

func (stopResults) isStopResults() {}

// FollowUp blocks completion and feeds reason back to the agent.
func (stopResults) FollowUp(reason string) StopOutput {
	return stopOutput{reason: reason}
}

// Noop returns an empty response (silent stdout).
func (stopResults) Noop() StopOutput {
	return stopOutput{}
}

func (stopOutput) allowedEvents() []string {
	return []string{EventAgentStop, EventSubagentStop}
}

func (o stopOutput) encode() ([]byte, int, error) {
	if o.reason == "" {
		return nil, 0, nil
	}
	out := map[string]any{
		"decision": "block",
		"reason":   o.reason,
	}
	b, err := json.Marshal(out)
	return b, 0, err
}
