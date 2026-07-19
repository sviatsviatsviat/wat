package agnostic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/cursor"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func resetTest(t *testing.T) {
	t.Helper()
	run.Reset()
}

func TestServe_MergeDispatch(t *testing.T) {
	resetTest(t)
	payload := `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`
	OnPostTool(func(ctx context.Context, hook PostToolHook, r PostToolResults) (PostToolResult, error) {
		return r.Context("first"), nil
	}).OnPostTool(func(ctx context.Context, hook PostToolHook, r PostToolResults) (PostToolResult, error) {
		return r.Context("second"), nil
	})

	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(payload), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	var got struct {
		HSO struct {
			AdditionalContext string `json:"additionalContext"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.HSO.AdditionalContext != "first\n\nsecond" {
		t.Fatalf("context = %q, want %q", got.HSO.AdditionalContext, "first\n\nsecond")
	}
}

func TestServe_DecisionPrecedence(t *testing.T) {
	resetTest(t)
	OnPreTool(func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		return r.Allow(), nil
	}).OnPreTool(func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		return r.Deny("blocked"), nil
	})

	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "deny") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestServe_WithDialectOverride(t *testing.T) {
	resetTest(t)
	OnPreTool(func(ctx context.Context, hook PreToolHook, _ PreToolResults) (PreToolResult, error) {
		return nil, nil
	})

	var stderr bytes.Buffer
	code := run.Serve(
		context.Background(),
		strings.NewReader(`{"sessionId":"s1","timestamp":1,"cwd":"/w"}`),
		new(bytes.Buffer),
		&stderr,
		run.WithDialect(copilot.Dialect),
	)
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "decode") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_CopilotPreTool(t *testing.T) {
	resetTest(t)
	OnPreTool(func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		return r.Deny("nope"), nil
	})

	var stdout, stderr bytes.Buffer
	code := run.Serve(
		context.Background(),
		strings.NewReader(copilotPreToolUse),
		&stdout,
		&stderr,
	)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "deny") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestServe_WithGetenv(t *testing.T) {
	resetTest(t)

	var stderr bytes.Buffer
	code := run.Serve(
		context.Background(),
		strings.NewReader(`{}`),
		new(bytes.Buffer),
		&stderr,
		run.WithGetenv(func(string) string { return "" }),
	)
	if code != 1 {
		t.Fatalf("exit = %d, want unknown dialect without env", code)
	}
	if !strings.Contains(stderr.String(), "unknown dialect") {
		t.Fatalf("stderr = %q", stderr.String())
	}

	stderr.Reset()
	code = run.Serve(
		context.Background(),
		strings.NewReader(`{"conversation_id":"c1","hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`),
		new(bytes.Buffer),
		&stderr,
		run.WithGetenv(func(key string) string {
			if key == "CURSOR_VERSION" {
				return "9.9.9"
			}
			return ""
		}),
	)
	if code != 0 {
		t.Fatalf("exit = %d, want 0 for env-detected cursor payload", code)
	}
}

func TestServe_HandlerErrorCopilotPreTool(t *testing.T) {
	resetTest(t)
	OnPreTool(func(ctx context.Context, hook PreToolHook, _ PreToolResults) (PreToolResult, error) {
		return nil, errors.New("boom")
	})

	var stderr bytes.Buffer
	code := run.Serve(
		context.Background(),
		strings.NewReader(copilotPreToolUse),
		new(bytes.Buffer),
		&stderr,
	)
	if code != copilot.HandlerErrorExit {
		t.Fatalf("exit = %d, want %d", code, copilot.HandlerErrorExit)
	}
	if !strings.Contains(stderr.String(), "boom") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_ZeroHandlers(t *testing.T) {
	resetTest(t)
	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q", stdout.Bytes())
	}
}

func TestServe_EmptyStdin(t *testing.T) {
	resetTest(t)
	var stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(""), new(bytes.Buffer), &stderr)
	if code != 1 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stderr.String(), "empty stdin") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_PreToolDenyAllAgents(t *testing.T) {
	denyShell := func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		if hook.Tool != nil && hook.Tool.Shell != "" {
			return r.Deny("destructive command blocked"), nil
		}
		return nil, nil
	}

	tests := []struct {
		name       string
		payload    string
		opts       []run.Option
		wantExit   int
		wantStdout func(string) bool
	}{
		{
			name:     "claude",
			payload:  claudePreToolUse,
			wantExit: 0,
			wantStdout: func(s string) bool {
				return strings.Contains(s, "deny") && strings.Contains(s, "destructive command blocked")
			},
		},
		{
			name:     "copilot",
			payload:  copilotPreToolUse,
			wantExit: 0,
			wantStdout: func(s string) bool {
				return strings.Contains(s, `"permission_decision":"deny"`) &&
					strings.Contains(s, "destructive command blocked")
			},
		},
		{
			name:     "cursor",
			payload:  cursorShell,
			wantExit: cursor.PermissionDenyExit,
			wantStdout: func(s string) bool {
				return strings.Contains(s, `"permission":"deny"`) &&
					strings.Contains(s, "destructive command blocked")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetTest(t)
			OnPreTool(denyShell)

			var stdout, stderr bytes.Buffer
			code := run.Serve(context.Background(), strings.NewReader(tt.payload), &stdout, &stderr, tt.opts...)
			if code != tt.wantExit {
				t.Fatalf("exit = %d, want %d; stderr = %q", code, tt.wantExit, stderr.String())
			}
			if !tt.wantStdout(stdout.String()) {
				t.Fatalf("stdout = %s", stdout.Bytes())
			}
		})
	}
}

func TestServe_AgnosticAndClaudeMerge(t *testing.T) {
	resetTest(t)

	payload := `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`
	OnPostTool(func(ctx context.Context, hook PostToolHook, r PostToolResults) (PostToolResult, error) {
		return r.Context("from-agnostic"), nil
	})
	claude.OnPostToolUse(func(ctx context.Context, hook run.Hook[claude.PostToolUse], r claude.PostToolUseResults) (claude.PostToolUseOutput, error) {
		return r.Context("from-claude"), nil
	})

	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(payload), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stdout.String(), "from-agnostic") || !strings.Contains(stdout.String(), "from-claude") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

type fakePreToolResult struct{}

func (fakePreToolResult) IsZero() bool { return false }
func (fakePreToolResult) WithUpdatedInput(map[string]any) PreToolResult {
	return fakePreToolResult{}
}

func TestServe_WrongPreToolResultType_FailOpen(t *testing.T) {
	resetTest(t)
	OnPreTool(func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		return fakePreToolResult{}, nil
	})

	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
	if code != 1 {
		t.Fatalf("exit = %d, want 1 (fail-open handler error)", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout should be empty on assert failure; got %s", stdout.Bytes())
	}
	if !strings.Contains(stderr.String(), "injected Results builder") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_CursorPreToolFanOut(t *testing.T) {
	denyAll := func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		return r.Deny("blocked"), nil
	}
	tests := []struct {
		name    string
		payload string
	}{
		{name: "beforeShell", payload: cursorShell},
		{name: "preToolUse", payload: cursorPreToolUse},
		{name: "beforeRead", payload: cursorBeforeRead},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetTest(t)
			OnPreTool(denyAll)
			var stdout, stderr bytes.Buffer
			code := run.Serve(context.Background(), strings.NewReader(tt.payload), &stdout, &stderr)
			if code != cursor.PermissionDenyExit {
				t.Fatalf("exit = %d, want %d; stderr = %q", code, cursor.PermissionDenyExit, stderr.String())
			}
			if !strings.Contains(stdout.String(), `"permission":"deny"`) {
				t.Fatalf("stdout = %s", stdout.Bytes())
			}
		})
	}
}

func TestServe_CursorAfterShellContext(t *testing.T) {
	resetTest(t)
	OnPostTool(func(ctx context.Context, hook PostToolHook, r PostToolResults) (PostToolResult, error) {
		return r.Context("after-shell note"), nil
	})
	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(cursorAfterShell), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "after-shell note") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestServe_CursorWithUpdatedInput_BeforeShellOmitsField(t *testing.T) {
	resetTest(t)
	OnPreTool(func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		return r.Allow().WithUpdatedInput(map[string]any{"command": "echo safe"}), nil
	})
	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(cursorShell), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d; stderr = %q", code, stderr.String())
	}
	out := stdout.String()
	if strings.Contains(out, "updated_input") {
		t.Fatalf("beforeShell should omit updated_input; stdout = %s", stdout.Bytes())
	}
	if !strings.Contains(out, `"permission":"allow"`) {
		t.Fatalf("beforeShell should retain allow; stdout = %s", stdout.Bytes())
	}
}

func TestServe_StopFollowUp(t *testing.T) {
	resetTest(t)
	OnStop(func(ctx context.Context, hook StopHook, r StopResults) (StopResult, error) {
		return r.FollowUp("keep going"), nil
	})
	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(claudeStop), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "keep going") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestServe_CopilotSubagentStopViaScopedStop(t *testing.T) {
	resetTest(t)
	var sawStop, sawSubagent bool
	OnStop(func(ctx context.Context, hook StopHook, r StopResults) (StopResult, error) {
		sawStop = true
		return r.FollowUp("keep-agent"), nil
	})
	OnSubagentStop(func(ctx context.Context, hook StopHook, r StopResults) (StopResult, error) {
		if hook.Subagent == nil || hook.Subagent.Type != "task" {
			t.Fatalf("Subagent = %+v", hook.Subagent)
		}
		sawSubagent = true
		return r.FollowUp("keep-subagent"), nil
	})
	payload := `{
  "hook_event_name": "Stop",
  "session_id": "s2",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "agent_name": "task",
  "stop_reason": "end_turn"
}`
	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(payload), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d; stderr = %q", code, stderr.String())
	}
	if sawStop {
		t.Fatal("OnStop should not run for subagent-scoped Stop")
	}
	if !sawSubagent {
		t.Fatal("OnSubagentStop did not run")
	}
	if !strings.Contains(stdout.String(), "keep-subagent") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
	if strings.Contains(stdout.String(), "keep-agent") {
		t.Fatalf("stdout should not include keep-agent: %s", stdout.Bytes())
	}
}

func TestServe_SessionStartContext(t *testing.T) {
	resetTest(t)
	OnSessionStart(func(ctx context.Context, hook SessionStartHook, r SessionStartResults) (SessionStartResult, error) {
		return r.Context("boot note"), nil
	})
	var stdout, stderr bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(claudeSessionStart), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d; stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "boot note") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}
