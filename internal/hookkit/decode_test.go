package hookkit

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestDecodeAsAndThen(t *testing.T) {
	t.Parallel()
	type ev struct {
		Name  string `json:"name"`
		Bound string
	}
	got, err := DecodeAsAndThen(json.RawMessage(`{"name":"x"}`), func(e *ev, _ []byte) {
		e.Bound = e.Name
	})
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "x" || got.Bound != "x" {
		t.Fatalf("got = %+v", got)
	}
	_, err = DecodeAsAndThen[ev]([]byte(`{`), nil)
	if err == nil {
		t.Fatal("expected unmarshal error")
	}
}

type wireEvent struct {
	HookEventName string `json:"hook_event_name"`
}

func (e wireEvent) EventName() string { return e.HookEventName }

func TestEventDecoder(t *testing.T) {
	t.Parallel()
	decodeErr := errors.New("decode")
	c := NewCodec("test", errors.New("empty"), decodeErr, errors.New("name"))
	c.Register("X", EventDecoder[wireEvent](c))

	ev, err := c.Decode([]byte(`{"hook_event_name":"X"}`))
	if err != nil {
		t.Fatal(err)
	}
	got, ok := ev.(wireEvent)
	if !ok || got.HookEventName != "X" {
		t.Fatalf("got = %#v", ev)
	}

	_, err = DecodeEvent[wireEvent](c, []byte(`{`), nil)
	if err == nil || !errors.Is(err, decodeErr) {
		t.Fatalf("err = %v", err)
	}
}
