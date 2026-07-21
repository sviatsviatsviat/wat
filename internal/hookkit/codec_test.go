package hookkit

import (
	"errors"
	"strings"
	"testing"
)

type namedEvent string

func (e namedEvent) EventName() string { return string(e) }

func TestCodec(t *testing.T) {
	t.Parallel()
	empty := errors.New("empty")
	decodeErr := errors.New("decode")
	nameRequired := errors.New("name required")

	c := NewCodec("test", empty, decodeErr, nameRequired, nil)
	c.Register("Known", func(raw []byte) (Event, error) {
		return namedEvent("Known"), nil
	})

	got, err := c.Decode([]byte(`{"hook_event_name":"Known"}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.EventName() != "Known" {
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

	_, err = c.Decode([]byte(`{"hook_event_name":"Other"}`))
	if err == nil {
		t.Fatal("expected error for unknown event")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestCodec_IsolatedRegistries(t *testing.T) {
	t.Parallel()
	empty := errors.New("empty")
	decodeErr := errors.New("decode")
	nameRequired := errors.New("name required")

	a := NewCodec("a", empty, decodeErr, nameRequired, nil)
	b := NewCodec("b", empty, decodeErr, nameRequired, nil)
	a.Register("X", func([]byte) (Event, error) { return namedEvent("a"), nil })

	_, err := b.Decode([]byte(`{"hook_event_name":"X"}`))
	if err == nil {
		t.Fatal("expected error: b must not see a's decoder")
	}
	if !strings.Contains(err.Error(), "unknown hook event") {
		t.Fatalf("error = %v", err)
	}
}

func TestCodec_WrapDecodeError(t *testing.T) {
	t.Parallel()
	decodeErr := errors.New("decode payload")
	c := NewCodec("claude", errors.New("empty"), decodeErr, errors.New("name"), nil)
	err := c.WrapDecodeError(struct{}{}, errors.New("bad json"))
	if !errors.Is(err, decodeErr) {
		t.Fatalf("errors.Is = false: %v", err)
	}
	if !strings.Contains(err.Error(), "claude: decode") {
		t.Fatalf("error = %v", err)
	}
}
