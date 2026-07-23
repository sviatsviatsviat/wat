package runtime

import (
	"errors"
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"testing"
)

func TestEventNameFromRaw(t *testing.T) {
	name, err := hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired).EventName([]byte(`{"hook_event_name":"PreToolUse","session_id":"s"}`))
	if err != nil {
		t.Fatal(err)
	}
	if name != "PreToolUse" {
		t.Fatalf("name = %q", name)
	}
	_, err = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired).EventName([]byte(`{"session_id":"s"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}
