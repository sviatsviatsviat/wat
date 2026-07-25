package runtime

import (
	"errors"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

func testCodec() *hookkit.Codec {
	return hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)
}

func TestDecode_UnknownEvent(t *testing.T) {
	_, err := testCodec().Decode([]byte(`{"session_id":"s1","hook_event_name":"FutureEvent","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDecode_RequiresHookEventName(t *testing.T) {
	_, err := testCodec().Decode([]byte(`{"session_id":"s1","cwd":"/w"}`))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	_, err := testCodec().Decode([]byte("not json"))
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrDecodePayload) {
		t.Fatalf("errors.Is ErrDecodePayload = false, err = %v", err)
	}
}

func TestEventNameFromRaw(t *testing.T) {
	c := testCodec()
	name, err := c.EventName([]byte(`{"hook_event_name":"PreToolUse","session_id":"s"}`))
	if err != nil {
		t.Fatal(err)
	}
	if name != "PreToolUse" {
		t.Fatalf("name = %q", name)
	}
	_, err = c.EventName([]byte(`{"session_id":"s"}`))
	if err == nil {
		t.Fatal("expected error without hook_event_name")
	}
	if !errors.Is(err, ErrEventNameRequired) {
		t.Fatalf("errors.Is ErrEventNameRequired = false, err = %v", err)
	}
}
