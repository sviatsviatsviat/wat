package hookkit

import (
	"encoding/json"
	"testing"
)

func TestCloneBytes(t *testing.T) {
	t.Parallel()
	raw := json.RawMessage(`{"a":1}`)
	got := CloneBytes(raw)
	if string(got) != string(raw) {
		t.Fatalf("CloneBytes() = %q, want %q", got, raw)
	}
	got[0] = '['
	if raw[0] == '[' {
		t.Fatal("CloneBytes must not alias input")
	}
	if CloneBytes(nil) != nil {
		t.Fatal("CloneBytes(nil) should be nil")
	}
}

func TestNullToNil(t *testing.T) {
	t.Parallel()
	if NullToNil(nil) != nil {
		t.Fatal("nil should stay nil")
	}
	if NullToNil(json.RawMessage("null")) != nil {
		t.Fatal("JSON null should become nil")
	}
	if NullToNil(json.RawMessage("  null\n")) != nil {
		t.Fatal("whitespace-padded JSON null should become nil")
	}
	in := json.RawMessage(`{"a":1}`)
	if string(NullToNil(in)) != `{"a":1}` {
		t.Fatalf("NullToNil = %s", NullToNil(in))
	}
	padded := json.RawMessage("  {\"a\":1}  ")
	if string(NullToNil(padded)) != string(padded) {
		t.Fatalf("non-null raw should be unchanged: %s", NullToNil(padded))
	}
}

func TestRawObjectField(t *testing.T) {
	t.Parallel()
	raw := json.RawMessage(`{"tool_input":{"command":"ls"},"other":null,"duration":0}`)
	got := RawObjectField(raw, "tool_input")
	if string(got) != `{"command":"ls"}` {
		t.Fatalf("RawObjectField = %q", got)
	}
	if RawObjectField(raw, "missing") != nil {
		t.Fatal("missing field should be nil")
	}
	if RawObjectField(raw, "other") != nil {
		t.Fatal("null field should be nil")
	}
	if got := RawObjectField(raw, "duration"); string(got) != "0" {
		t.Fatalf("explicit zero field = %q, want 0", got)
	}
}

func TestRawToText(t *testing.T) {
	t.Parallel()
	if got := RawToText([]byte(`"hello"`)); got != "hello" {
		t.Fatalf("string = %q", got)
	}
	if got := RawToText([]byte(`{"k":"v"}`)); got != `{"k":"v"}` {
		t.Fatalf("object = %q", got)
	}
	if got := RawToText(nil); got != "" {
		t.Fatalf("empty = %q", got)
	}
}
