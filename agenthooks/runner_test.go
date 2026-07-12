package agenthooks

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestMux_Serve_MergeDispatch(t *testing.T) {
	mux := NewMux()
	mux.OnAny(func(ctx context.Context, ev *Event) (Result, error) {
		return Context("from-any"), nil
	})
	mux.On(KindPreTool, func(ctx context.Context, ev *Event) (Result, error) {
		return Context("from-kind"), nil
	})

	var stdout, stderr bytes.Buffer
	code := mux.Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if stdout.Len() == 0 {
		t.Fatal("expected stdout")
	}
	var got struct {
		HSO struct {
			AdditionalContext string `json:"additionalContext"`
		} `json:"hookSpecificOutput"`
	}
	if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	want := "from-any\n\nfrom-kind"
	if got.HSO.AdditionalContext != want {
		t.Fatalf("context = %q, want %q", got.HSO.AdditionalContext, want)
	}
}

func TestMux_Serve_DecisionPrecedence(t *testing.T) {
	mux := NewMux()
	mux.OnAny(func(ctx context.Context, ev *Event) (Result, error) {
		return Allow(), nil
	})
	mux.On(KindPreTool, func(ctx context.Context, ev *Event) (Result, error) {
		return Deny("blocked"), nil
	})

	var stdout, stderr bytes.Buffer
	code := mux.Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "deny") {
		t.Fatalf("stdout = %s", stdout.Bytes())
	}
}

func TestMux_Serve_WithDialectOverride(t *testing.T) {
	mux := NewMux()
	mux.On(KindPreTool, func(ctx context.Context, ev *Event) (Result, error) {
		return Result{}, nil
	})

	var stderr bytes.Buffer
	code := mux.Serve(
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

func TestMux_Serve_WithEvent(t *testing.T) {
	mux := NewMux()
	mux.On(KindPreTool, func(ctx context.Context, ev *Event) (Result, error) {
		return Deny("nope"), nil
	})

	var stdout, stderr bytes.Buffer
	code := mux.Serve(
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

func TestMux_Serve_WithGetenv(t *testing.T) {
	mux := NewMux()

	var stderr bytes.Buffer
	code := mux.Serve(
		context.Background(),
		strings.NewReader(`{}`),
		new(bytes.Buffer),
		&stderr,
	)
	if code != 1 {
		t.Fatalf("exit = %d, want unknown dialect without env", code)
	}
	if !strings.Contains(stderr.String(), "unknown dialect") {
		t.Fatalf("stderr = %q", stderr.String())
	}

	stderr.Reset()
	code = mux.Serve(
		context.Background(),
		strings.NewReader(`{}`),
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

func TestMux_Serve_HandlerErrorCopilotPreTool(t *testing.T) {
	mux := NewMux()
	mux.On(KindPreTool, func(ctx context.Context, ev *Event) (Result, error) {
		return Result{}, errors.New("boom")
	})

	var stderr bytes.Buffer
	code := mux.Serve(
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

func TestMux_Serve_ZeroHandlers(t *testing.T) {
	mux := NewMux()
	var stdout, stderr bytes.Buffer
	code := mux.Serve(context.Background(), strings.NewReader(claudePreToolUse), &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q", stdout.Bytes())
	}
}

func TestMux_Serve_EmptyStdin(t *testing.T) {
	mux := NewMux()
	var stderr bytes.Buffer
	code := mux.Serve(context.Background(), strings.NewReader(""), new(bytes.Buffer), &stderr)
	if code != 1 {
		t.Fatalf("exit = %d", code)
	}
	if !strings.Contains(stderr.String(), "empty stdin") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestMux_Serve_PreToolDenyAllAgents(t *testing.T) {
	denyShell := func(ctx context.Context, ev *Event) (Result, error) {
		if ev.Tool != nil && ev.Tool.Shell != "" {
			return Deny("destructive command blocked"), nil
		}
		return Result{}, nil
	}

	tests := []struct {
		name       string
		payload    string
		opts       []Option
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
			opts:     []Option{WithEvent("preToolUse")},
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
			mux := NewMux()
			mux.On(KindPreTool, denyShell)

			var stdout, stderr bytes.Buffer
			code := mux.Serve(context.Background(), strings.NewReader(tt.payload), &stdout, &stderr, tt.opts...)
			if code != tt.wantExit {
				t.Fatalf("exit = %d, want %d; stderr = %q", code, tt.wantExit, stderr.String())
			}
			if !tt.wantStdout(stdout.String()) {
				t.Fatalf("stdout = %s", stdout.Bytes())
			}
		})
	}
}

func TestMux_OnNilIgnored(t *testing.T) {
	mux := NewMux()
	mux.On(KindPreTool, nil)
	mux.OnAny(nil)
	if len(mux.kindHandlers[KindPreTool]) != 0 {
		t.Fatal("nil On handler should be ignored")
	}
	if len(mux.anyHandlers) != 0 {
		t.Fatal("nil OnAny handler should be ignored")
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
			if got := handlerErrorExit(tt.ev, errors.New("x")); got != tt.want {
				t.Fatalf("handlerErrorExit() = %d, want %d", got, tt.want)
			}
		})
	}
}
