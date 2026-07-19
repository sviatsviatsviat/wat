package hookkit

import (
	"errors"
	"strings"
	"testing"
)

func TestCodec(t *testing.T) {
	t.Parallel()
	empty := errors.New("empty")
	decodeErr := errors.New("decode")
	nameRequired := errors.New("name required")

	c := NewCodec("test", empty, decodeErr, nameRequired)
	c.Register("Known", func(raw []byte) (any, error) {
		return string(raw), nil
	})

	got, err := c.Decode([]byte(`{"hook_event_name":"Known"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got != `{"hook_event_name":"Known"}` {
		t.Fatalf("got = %v", got)
	}

	name, err := c.EventName([]byte(`{"hook_event_name":"Known"}`))
	if err != nil || name != "Known" {
		t.Fatalf("EventName = %q, %v", name, err)
	}

	_, err = c.Decode(nil)
	if !errors.Is(err, empty) {
		t.Fatalf("empty: %v", err)
	}

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic for unknown event")
		}
		msg, ok := r.(string)
		if !ok || !strings.Contains(msg, "unknown hook event") {
			t.Fatalf("recover = %#v", r)
		}
	}()
	_, _ = c.Decode([]byte(`{"hook_event_name":"Other"}`))
}

func TestCodec_IsolatedRegistries(t *testing.T) {
	t.Parallel()
	empty := errors.New("empty")
	decodeErr := errors.New("decode")
	nameRequired := errors.New("name required")

	a := NewCodec("a", empty, decodeErr, nameRequired)
	b := NewCodec("b", empty, decodeErr, nameRequired)
	a.Register("X", func([]byte) (any, error) { return "a", nil })

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic: b must not see a's decoder")
		}
	}()
	_, _ = b.Decode([]byte(`{"hook_event_name":"X"}`))
}

func TestCodec_WrapDecodeError(t *testing.T) {
	t.Parallel()
	decodeErr := errors.New("decode payload")
	c := NewCodec("claude", errors.New("empty"), decodeErr, errors.New("name"))
	err := c.WrapDecodeError(struct{}{}, errors.New("bad json"))
	if !errors.Is(err, decodeErr) {
		t.Fatalf("errors.Is = false: %v", err)
	}
	if !strings.Contains(err.Error(), "claude: decode") {
		t.Fatalf("error = %v", err)
	}
}
