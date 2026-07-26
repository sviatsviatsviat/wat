package sessionstart

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/event"
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

func TestEncode_SessionStartEnv(t *testing.T) {
	out, code, err := results{}.Noop().
		WithEnv(map[string]string{"K": "V"}).
		WithAdditionalContext("ctx").Encode()
	if err != nil || code != 0 {
		t.Fatalf("encode: %v code=%d", err, code)
	}
	var got struct {
		Env      map[string]string `json:"env"`
		Ctx      string            `json:"additional_context"`
		Continue *bool             `json:"continue"`
		UserMsg  string            `json:"user_message"`
	}
	if err := json.Unmarshal(out, &got); err != nil {
		t.Fatal(err)
	}
	if got.Env["K"] != "V" || got.Ctx != "ctx" {
		t.Fatalf("bad output: %s", out)
	}
	if got.Continue != nil || got.UserMsg != "" {
		t.Fatalf("encode must not emit continue/user_message: %s", out)
	}
	if strings.Contains(string(out), "continue") || strings.Contains(string(out), "user_message") {
		t.Fatalf("encode must not emit continue/user_message keys: %s", out)
	}
}

func TestDecode_SessionStart(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want Event
	}{
		{
			name: "background without composer_mode",
			raw:  `{"hook_event_name":"sessionStart","conversation_id":"c1","model":"gpt","is_background_agent":true,"cwd":"/w"}`,
			want: Event{
				Envelope:          eventEnvelope("c1", "gpt", "/w"),
				IsBackgroundAgent: true,
			},
		},
		{
			name: "composer_mode agent",
			raw:  `{"hook_event_name":"sessionStart","conversation_id":"c1","model":"gpt","is_background_agent":false,"composer_mode":"agent","cwd":"/w"}`,
			want: Event{
				Envelope:          eventEnvelope("c1", "gpt", "/w"),
				IsBackgroundAgent: false,
				ComposerMode:      "agent",
			},
		},
		{
			name: "composer_mode ask",
			raw:  `{"hook_event_name":"sessionStart","conversation_id":"c2","model":"gpt","is_background_agent":false,"composer_mode":"ask","cwd":"/w"}`,
			want: Event{
				Envelope:     eventEnvelope("c2", "gpt", "/w"),
				ComposerMode: "ask",
			},
		},
		{
			name: "composer_mode edit",
			raw:  `{"hook_event_name":"sessionStart","conversation_id":"c3","model":"gpt","is_background_agent":false,"composer_mode":"edit","cwd":"/w"}`,
			want: Event{
				Envelope:     eventEnvelope("c3", "gpt", "/w"),
				ComposerMode: "edit",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mustDecode[Event](t, tt.raw)
			if got.ConversationID != tt.want.ConversationID ||
				got.Model != tt.want.Model ||
				got.Cwd != tt.want.Cwd ||
				got.IsBackgroundAgent != tt.want.IsBackgroundAgent ||
				got.ComposerMode != tt.want.ComposerMode {
				t.Fatalf("event=%+v want=%+v", got, tt.want)
			}
		})
	}
}

func eventEnvelope(conversationID, model, cwd string) event.Envelope {
	return event.Envelope{
		ConversationID: conversationID,
		Model:          model,
		HookEventName:  "sessionStart",
		Cwd:            cwd,
	}
}

func TestMerge_SessionStart_contextJoins(t *testing.T) {
	a := results{}.Context("one")
	b := results{}.Context("two")
	merged, warnings, err := a.Merge(b.(hookkit.Output))
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	out := merged.(output)
	if out.additionalContext != "one\n\ntwo" {
		t.Fatalf("context = %q", out.additionalContext)
	}
	if merged.Stop() {
		t.Fatal("context should not stop")
	}
}

func init() {
	register(testCodec)
}
