package posttooluse

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"

	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

var testCodec = hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)

func mustDecode[E any](t *testing.T, raw string) E {
	t.Helper()
	ev, err := testCodec.Decode([]byte(raw))
	if err != nil {
		t.Fatal(err)
	}
	if ev.EventName() == "" {
		t.Fatal("EventName empty")
	}
	typed, ok := ev.(E)
	if !ok {
		t.Fatalf("want %T, got %T", *new(E), ev)
	}
	return typed
}

func TestDecode_PostToolUse(t *testing.T) {
	const raw = `{
		"hook_event_name":"postToolUse",
		"conversation_id":"c1",
		"generation_id":"g1",
		"model":"claude-opus-4-7-thinking-max",
		"model_id":"claude-opus-4-7",
		"model_params":[
			{"id":"thinking","value":"true"},
			{"id":"context","value":"1m"},
			{"id":"effort","value":"max"}
		],
		"cwd":"/project",
		"tool_name":"Shell",
		"tool_use_id":"abc123",
		"tool_input":{"command":"npm test"},
		"tool_output":"{\"exitCode\":0,\"stdout\":\"All tests passed\"}",
		"duration":5432
	}`
	e := mustDecode[Event](t, raw)
	if e.ToolName != "Shell" || e.ToolUseID != "abc123" {
		t.Fatalf("tool identity: name=%q id=%q", e.ToolName, e.ToolUseID)
	}
	if e.ToolOutput != `{"exitCode":0,"stdout":"All tests passed"}` {
		t.Fatalf("ToolOutput=%q", e.ToolOutput)
	}
	if e.DurationMillis() != 5432 {
		t.Fatalf("DurationMillis()=%d, want 5432", e.DurationMillis())
	}
	if e.Model != "claude-opus-4-7-thinking-max" || e.ModelID != "claude-opus-4-7" {
		t.Fatalf("model=%q model_id=%q", e.Model, e.ModelID)
	}
	if len(e.ModelParams) != 3 {
		t.Fatalf("ModelParams len=%d, want 3", len(e.ModelParams))
	}
	if e.ModelParams[0].ID != "thinking" || e.ModelParams[0].Value != "true" {
		t.Fatalf("ModelParams[0]=%+v", e.ModelParams[0])
	}
	if e.Cwd != "/project" || e.GenerationID != "g1" {
		t.Fatalf("cwd=%q generation_id=%q", e.Cwd, e.GenerationID)
	}
	cmd, ok := e.ToolInput.AsShell()
	if !ok || cmd.Command != "npm test" {
		t.Fatalf("ToolInput AsShell=%v ok=%v", cmd, ok)
	}
}

func TestDurationMillis_PrefersDurationOverDurationMs(t *testing.T) {
	e := mustDecode[Event](t, `{
		"hook_event_name":"postToolUse",
		"conversation_id":"c1",
		"tool_name":"Read",
		"tool_output":"contents",
		"duration":100,
		"duration_ms":999
	}`)
	if e.DurationMillis() != 100 {
		t.Fatalf("DurationMillis()=%d, want documented duration field 100", e.DurationMillis())
	}
}

func TestDurationMillis_FallsBackToDurationMs(t *testing.T) {
	e := mustDecode[Event](t, `{
		"hook_event_name":"postToolUse",
		"conversation_id":"c1",
		"tool_name":"Read",
		"tool_output":"contents",
		"duration_ms":50
	}`)
	if e.DurationMillis() != 50 {
		t.Fatalf("DurationMillis()=%d, want duration_ms fallback 50", e.DurationMillis())
	}
}

func init() {
	register(testCodec)
}
