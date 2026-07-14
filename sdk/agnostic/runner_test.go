package agnostic

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/claude"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func resetTest(t *testing.T) {
	t.Helper()
	ResetHandlers()
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
	code := Serve(context.Background(), strings.NewReader(payload), &stdout, &stderr)
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
	code := Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
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
		return PreToolResult{}, nil
	})

	var stderr bytes.Buffer
	code := Serve(
		context.Background(),
		strings.NewReader(copilotCamelPreToolUse),
		new(bytes.Buffer),
		&stderr,
		WithDialect(Copilot),
	)
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "decode") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_WithEvent(t *testing.T) {
	resetTest(t)
	OnPreTool(func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		return r.Deny("nope"), nil
	})

	var stdout, stderr bytes.Buffer
	code := Serve(
		context.Background(),
		strings.NewReader(copilotCamelPreToolUse),
		&stdout,
		&stderr,
		WithEvent("preToolUse"),
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
	code := Serve(
		context.Background(),
		strings.NewReader(`{}`),
		new(bytes.Buffer),
		&stderr,
		WithGetenv(func(string) string { return "" }),
	)
	if code != 1 {
		t.Fatalf("exit = %d, want unknown dialect without env", code)
	}
	if !strings.Contains(stderr.String(), "unknown dialect") {
		t.Fatalf("stderr = %q", stderr.String())
	}

	stderr.Reset()
	code = Serve(
		context.Background(),
		strings.NewReader(`{"conversation_id":"c1","hook_event_name":"workspaceOpen","cursor_version":"1.7.2","workspace_roots":["/w"]}`),
		new(bytes.Buffer),
		&stderr,
		WithGetenv(func(key string) string {
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
		return PreToolResult{}, errors.New("boom")
	})

	var stderr bytes.Buffer
	code := Serve(
		context.Background(),
		strings.NewReader(copilotCamelPreToolUse),
		new(bytes.Buffer),
		&stderr,
		WithEvent("preToolUse"),
	)
	if code != CopilotPreToolErrorExit {
		t.Fatalf("exit = %d, want %d", code, CopilotPreToolErrorExit)
	}
	if !strings.Contains(stderr.String(), "boom") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_ZeroHandlers(t *testing.T) {
	resetTest(t)
	var stdout, stderr bytes.Buffer
	code := Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
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
	code := Serve(context.Background(), strings.NewReader(""), new(bytes.Buffer), &stderr)
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
		return PreToolResult{}, nil
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
			payload:  copilotCamelPreToolUse,
			opts:     []run.Option{WithEvent("preToolUse")},
			wantExit: 0,
			wantStdout: func(s string) bool {
				return strings.Contains(s, `"permissionDecision":"deny"`) &&
					strings.Contains(s, "destructive command blocked")
			},
		},
		{
			name:     "cursor",
			payload:  cursorShell,
			wantExit: CursorWarnExit,
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
			code := Serve(context.Background(), strings.NewReader(tt.payload), &stdout, &stderr, tt.opts...)
			if code != tt.wantExit {
				t.Fatalf("exit = %d, want %d; stderr = %q", code, tt.wantExit, stderr.String())
			}
			if !tt.wantStdout(stdout.String()) {
				t.Fatalf("stdout = %s", stdout.Bytes())
			}
		})
	}
}

func TestServe_OnAnyObserveOnly(t *testing.T) {
	resetTest(t)
	OnAny(func(ctx context.Context, hook AnyHook) error {
		return nil
	})
	OnPreTool(func(ctx context.Context, hook PreToolHook, r PreToolResults) (PreToolResult, error) {
		return r.Deny("blocked"), nil
	})

	var stdout, stderr bytes.Buffer
	code := Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "deny") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestResetHandlers_OwnerScoped(t *testing.T) {
	run.Reset()
	claude.ResetHandlers()

	payload := `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`
	OnPostTool(func(ctx context.Context, hook PostToolHook, r PostToolResults) (PostToolResult, error) {
		return r.Context("from-agnostic"), nil
	})
	new(claude.Chain).PostToolUse(func(ctx context.Context, hook claude.PostToolUseHook, r claude.PostToolUseResults) (claude.PostToolUseOutput, error) {
		return r.Context("from-claude"), nil
	})

	ResetHandlers()

	var stdout bytes.Buffer
	code := run.Serve(context.Background(), strings.NewReader(payload), &stdout, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if strings.Contains(stdout.String(), "from-agnostic") {
		t.Fatalf("agnostic handler should be cleared; stdout = %s", stdout.Bytes())
	}
	if !strings.Contains(stdout.String(), "from-claude") {
		t.Fatalf("claude handler should remain; stdout = %s", stdout.Bytes())
	}
}

func TestServe_AgnosticAndClaudeMerge(t *testing.T) {
	resetTest(t)
	claude.ResetHandlers()

	payload := `{"session_id":"s","hook_event_name":"PostToolUse","tool_name":"Read","tool_response":"file contents"}`
	OnPostTool(func(ctx context.Context, hook PostToolHook, r PostToolResults) (PostToolResult, error) {
		return r.Context("from-agnostic"), nil
	})
	new(claude.Chain).PostToolUse(func(ctx context.Context, hook claude.PostToolUseHook, r claude.PostToolUseResults) (claude.PostToolUseOutput, error) {
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

func TestHandlerErrorExit(t *testing.T) {
	tests := []struct {
		name string
		ev   *Event
		want int
	}{
		{
			name: "copilot pre tool",
			ev:   &Event{Agent: Copilot, Kind: KindPreTool},
			want: CopilotPreToolErrorExit,
		},
		{
			name: "copilot stop",
			ev:   &Event{Agent: Copilot, Kind: KindStop},
			want: 1,
		},
		{
			name: "cursor",
			ev:   &Event{Agent: Cursor, Kind: KindPreTool},
			want: CursorHandlerErrorExit,
		},
		{
			name: "claude",
			ev:   &Event{Agent: Claude, Kind: KindPreTool},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := handlerErrorExit(tt.ev); got != tt.want {
				t.Fatalf("handlerErrorExit() = %d, want %d", got, tt.want)
			}
		})
	}
}
