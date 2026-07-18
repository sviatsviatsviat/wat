package hookkit

import (
	"encoding/json"
	"testing"
)

type testEnvelope struct {
	raw json.RawMessage
}

func (e *testEnvelope) DecodedRaw() json.RawMessage {
	return e.raw
}

func TestRawBytes(t *testing.T) {
	t.Parallel()
	env := &testEnvelope{raw: json.RawMessage(`{"x":1}`)}
	got := RawBytes(env, env)
	if string(got) != `{"x":1}` {
		t.Fatalf("RawBytes() = %q", got)
	}

	got = RawBytes(map[string]int{"y": 2}, nil)
	if string(got) != `{"y":2}` {
		t.Fatalf("RawBytes(marshal fallback) = %q", got)
	}

	got = RawBytes(nil, nil)
	if got != nil {
		t.Fatalf("RawBytes(nil) = %q, want nil", got)
	}

	empty := &testEnvelope{}
	got = RawBytes(map[string]int{"z": 3}, empty)
	if string(got) != `{"z":3}` {
		t.Fatalf("RawBytes(empty accessor) = %q, want marshal fallback", got)
	}
}
