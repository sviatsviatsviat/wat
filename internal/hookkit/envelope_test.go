package hookkit

import (
	"encoding/json"
	"testing"
)

type testRaw struct {
	raw json.RawMessage
}

func (e *testRaw) Raw() json.RawMessage {
	return e.raw
}

func TestEventRaw(t *testing.T) {
	t.Parallel()
	acc := &testRaw{raw: json.RawMessage(`{"x":1}`)}
	got := EventRaw(acc, acc)
	if string(got) != `{"x":1}` {
		t.Fatalf("EventRaw() = %q", got)
	}

	got = EventRaw(map[string]int{"y": 2}, nil)
	if string(got) != `{"y":2}` {
		t.Fatalf("EventRaw(marshal fallback) = %q", got)
	}

	got = EventRaw(nil, nil)
	if got != nil {
		t.Fatalf("EventRaw(nil) = %q, want nil", got)
	}

	empty := &testRaw{}
	got = EventRaw(map[string]int{"z": 3}, empty)
	if string(got) != `{"z":3}` {
		t.Fatalf("EventRaw(empty accessor) = %q, want marshal fallback", got)
	}
}
