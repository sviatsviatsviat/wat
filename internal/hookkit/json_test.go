package hookkit

import (
	"encoding/json"
	"testing"
)

func TestCloneRaw(t *testing.T) {
	t.Parallel()
	raw := json.RawMessage(`{"a":1}`)
	got := CloneRaw(raw)
	if string(got) != string(raw) {
		t.Fatalf("CloneRaw() = %q, want %q", got, raw)
	}
	got[0] = '['
	if raw[0] == '[' {
		t.Fatal("CloneRaw must not alias input")
	}
	if CloneRaw(nil) != nil {
		t.Fatal("CloneRaw(nil) should be nil")
	}
}

func TestRawObjectField(t *testing.T) {
	t.Parallel()
	raw := json.RawMessage(`{"tool_input":{"command":"ls"},"other":null}`)
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
