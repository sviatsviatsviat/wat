package claude_test

import (
	"context"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestUseHooks_registersPreviouslyDecodeOnlyEvents(t *testing.T) {
	hooks := claude.UseHooks().
		Setup(func(context.Context, claude.Setup, claude.SetupResults) (claude.SetupOutput, error) {
			return nil, nil
		}).
		InstructionsLoaded(func(context.Context, claude.InstructionsLoaded) error { return nil }).
		PostToolBatch(func(context.Context, claude.PostToolBatch, claude.PostToolBatchResults) (claude.DecisionOutput, error) {
			return nil, nil
		}).
		TeammateIdle(func(context.Context, claude.TeammateIdle, claude.TeammateIdleResults) (claude.ExitBlockOutput, error) {
			return nil, nil
		}).
		ConfigChange(func(context.Context, claude.ConfigChange, claude.ConfigChangeResults) (claude.ConfigChangeOutput, error) {
			return nil, nil
		}).
		CwdChanged(func(context.Context, claude.CwdChanged, claude.CwdChangedResults) (claude.CwdChangedOutput, error) {
			return nil, nil
		}).
		FileChanged(func(context.Context, claude.FileChanged, claude.FileChangedResults) (claude.FileChangedOutput, error) {
			return nil, nil
		}).
		WorktreeRemove(func(context.Context, claude.WorktreeRemove) error { return nil }).
		PostCompact(func(context.Context, claude.PostCompact) error { return nil }).
		ElicitationResult(func(context.Context, claude.ElicitationResult, claude.ElicitationResultResults) (claude.ElicitationResultOutput, error) {
			return nil, nil
		}).
		StopFailure(func(context.Context, claude.StopFailure) error { return nil })

	got := map[string]int{}
	for _, e := range run.Inspect(hooks).Registrations {
		if e.Dialect != claude.Dialect {
			t.Fatalf("unexpected dialect %q", e.Dialect)
		}
		got[e.Event] = e.HandlerCount
	}
	want := []string{
		claude.EventSetup,
		claude.EventInstructionsLoaded,
		claude.EventPostToolBatch,
		claude.EventTeammateIdle,
		claude.EventConfigChange,
		claude.EventCwdChanged,
		claude.EventFileChanged,
		claude.EventWorktreeRemove,
		claude.EventPostCompact,
		claude.EventElicitationResult,
		claude.EventStopFailure,
	}
	for _, name := range want {
		if got[name] != 1 {
			t.Fatalf("Inspect[%q] = %d, want 1 (map=%v)", name, got[name], got)
		}
	}
}
