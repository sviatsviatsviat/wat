package agnostic_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestInspect_expandsPortableHooksToNativeRegistrations(t *testing.T) {
	portable := agnostic.UseHooks().
		OnPreTool(func(context.Context, agnostic.PreToolEvent, agnostic.PreToolResults) (agnostic.PreToolResult, error) {
			return nil, nil
		}).
		OnStop(func(context.Context, agnostic.StopEvent, agnostic.StopResults) (agnostic.StopResult, error) {
			return nil, nil
		})
	native := cursor.UseHooks().
		BeforeShellExecution(func(context.Context, cursor.BeforeShellExecution, cursor.PermissionResults) (cursor.PermissionOutput, error) {
			return nil, nil
		})

	manifest := run.Inspect(portable, native)
	want := map[string][]string{
		claude.Dialect:  {claude.EventPreToolUse, claude.EventStop},
		copilot.Dialect: {copilot.EventPreToolUse, copilot.EventAgentStop},
		cursor.Dialect: {
			cursor.EventBeforeMCPExecution,
			cursor.EventBeforeReadFile,
			cursor.EventBeforeShellExecution,
			cursor.EventPreToolUse,
			cursor.EventStop,
		},
	}
	for dialect, events := range want {
		if got := manifest.EventsFor(dialect); !reflect.DeepEqual(got, events) {
			t.Fatalf("EventsFor(%s) = %v, want %v", dialect, got, events)
		}
	}

	var beforeShellHandlers int
	for _, registration := range manifest.Registrations {
		if registration.Dialect == cursor.Dialect && registration.Event == cursor.EventBeforeShellExecution {
			beforeShellHandlers = registration.HandlerCount
		}
	}
	if beforeShellHandlers != 2 {
		t.Fatalf("beforeShellExecution handler count = %d, want 2", beforeShellHandlers)
	}
}

func TestInspect_OnPostToolIncludesObserveAfterMCPExecution(t *testing.T) {
	portable := agnostic.UseHooks().OnPostTool(func(context.Context, agnostic.PostToolEvent, agnostic.PostToolResults) (agnostic.PostToolResult, error) {
		return nil, nil
	})
	manifest := run.Inspect(portable)
	got := manifest.EventsFor(cursor.Dialect)
	want := []string{
		cursor.EventAfterFileEdit,
		cursor.EventAfterMCPExecution,
		cursor.EventAfterShellExecution,
		cursor.EventPostToolUse,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("EventsFor(cursor) = %v, want %v", got, want)
	}
}
