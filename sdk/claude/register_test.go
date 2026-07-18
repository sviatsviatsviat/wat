package claude_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestServe_PreToolDeny(t *testing.T) {
	run.Reset()
	claude.OnPreToolUse(func(ctx context.Context, hook claude.Hook[claude.PreToolUse], r claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		return r.Deny("destructive command"), nil
	})

	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(preToolUsePayload), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stdout.String(), `"permissionDecision":"deny"`) {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestServe_FailPolicy(t *testing.T) {
	run.Reset()
	claude.OnPreToolUse(func(ctx context.Context, hook claude.Hook[claude.PreToolUse], _ claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		return nil, context.Canceled
	})

	code := run.Serve(context.Background(), strings.NewReader(preToolUsePayload), &bytes.Buffer{}, &bytes.Buffer{}, claude.WithFailPolicy(claude.FailOpen))
	if code != claude.HandlerErrorExit {
		t.Fatalf("FailOpen exit = %d, want %d", code, claude.HandlerErrorExit)
	}
	run.Reset()
	claude.OnPreToolUse(func(ctx context.Context, hook claude.Hook[claude.PreToolUse], _ claude.PreToolUseResults) (claude.PreToolUseOutput, error) {
		return nil, context.Canceled
	})
	code = run.Serve(context.Background(), strings.NewReader(preToolUsePayload), &bytes.Buffer{}, &bytes.Buffer{}, claude.WithFailPolicy(claude.FailBlock))
	if code != claude.FailBlockExit {
		t.Fatalf("FailBlock exit = %d, want %d", code, claude.FailBlockExit)
	}
}
