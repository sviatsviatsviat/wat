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
