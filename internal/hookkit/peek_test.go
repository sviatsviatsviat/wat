package hookkit

import "testing"

func TestPeekHookEventName(t *testing.T) {
	t.Parallel()
	got, err := PeekHookEventName([]byte(`{"hook_event_name":"PreToolUse","cwd":"/w"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got != "PreToolUse" {
		t.Fatalf("PeekHookEventName = %q, want PreToolUse", got)
	}
	got, err = PeekHookEventName([]byte(`{"cwd":"/w"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Fatalf("missing name = %q, want empty", got)
	}
	_, err = PeekHookEventName([]byte(`{`))
	if err == nil {
		t.Fatal("expected JSON error")
	}
}
