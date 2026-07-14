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
	got := RawBytes(env, nil, env, nil)
	if string(got) != `{"x":1}` {
		t.Fatalf("RawBytes() = %q", got)
	}
	raw := json.RawMessage(`{"fallback":true}`)
	got = RawBytes(nil, raw, nil, nil)
	if string(got) != string(raw) {
		t.Fatalf("RawBytes(raw event) = %q", got)
	}
	empty := json.RawMessage{}
	got = RawBytes(nil, empty, env, nil)
	if string(got) != `{"x":1}` {
		t.Fatalf("RawBytes(empty raw) = %q, want accessor fallback", got)
	}
}
