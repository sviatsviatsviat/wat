package copilot

import (
	"testing"

	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

func TestMapAgentStop_StopHookActive(t *testing.T) {
	ev := mapAgentStop(sdkcopilot.AgentStop{
		StopReason:     "end_turn",
		StopHookActive: true,
	})
	if ev.Turn == nil || !ev.Turn.StopHookActive || ev.Turn.Status != "end_turn" {
		t.Fatalf("Turn=%+v", ev.Turn)
	}
	if ev.Subagent != nil {
		t.Fatalf("Subagent=%+v, want nil", ev.Subagent)
	}
}

func TestMapAgentStopAsSubagent_StopHookActive(t *testing.T) {
	hook := sdkcopilot.AgentStop{
		StopReason:     "end_turn",
		StopHookActive: true,
	}
	hook.AgentName = "task"
	ev := mapAgentStopAsSubagent(hook)
	if ev.Turn == nil || !ev.Turn.StopHookActive || ev.Turn.Status != "end_turn" {
		t.Fatalf("Turn=%+v", ev.Turn)
	}
	if ev.Subagent == nil || ev.Subagent.Type != "task" {
		t.Fatalf("Subagent=%+v", ev.Subagent)
	}
}
