package run

import (
	"bytes"
	"context"
	"strings"
	"sync/atomic"
	"testing"
)

type testCodec struct {
	eventName   string
	decodeCalls *atomic.Int32
}

type testEvent struct {
	raw string
}

func (testEvent) EventName() string { return "TestEvent" }

type emptyOutput struct{}

func (emptyOutput) IsZero() bool { return true }

func (emptyOutput) Encode() ([]byte, int, error) { return nil, 0, nil }

func (c testCodec) EventName([]byte) (string, error) {
	return c.eventName, nil
}

func (c testCodec) Decode(raw []byte) (Event, error) {
	if c.decodeCalls != nil {
		c.decodeCalls.Add(1)
	}
	return testEvent{raw: string(raw)}, nil
}

func TestServe_DecodesOnce(t *testing.T) {
	r := NewRegistry()
	var decodeCalls atomic.Int32
	r.registerDialect("testdialect", DialectOps{
		Detect: func([]byte, func(string) string) bool { return true },
		Codec:  testCodec{eventName: "TestEvent", decodeCalls: &decodeCalls},
	})
	var handlerCalls atomic.Int32
	for range 3 {
		r.RegisterHandler("testdialect", Handler(func(_ context.Context, hook Hook[testEvent]) (emptyOutput, error) {
			handlerCalls.Add(1)
			if hook.Event.raw != `{"ok":true}` {
				t.Errorf("event = %#v", hook.Event)
			}
			return emptyOutput{}, nil
		}))
	}

	code := r.serve(context.Background(), strings.NewReader(`{"ok":true}`), &bytes.Buffer{}, &bytes.Buffer{}, applyOptions())
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if got := decodeCalls.Load(); got != 1 {
		t.Fatalf("Decode calls = %d, want 1", got)
	}
	if got := handlerCalls.Load(); got != 3 {
		t.Fatalf("handler calls = %d, want 3", got)
	}
}

func TestServe_SkipsDecodeWhenNoHandlers(t *testing.T) {
	r := NewRegistry()
	var decodeCalls atomic.Int32
	r.registerDialect("empty", DialectOps{
		Detect: func([]byte, func(string) string) bool { return true },
		Codec:  testCodec{eventName: "NoHandlers", decodeCalls: &decodeCalls},
	})
	code := r.serve(context.Background(), strings.NewReader(`{}`), &bytes.Buffer{}, &bytes.Buffer{}, applyOptions())
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if got := decodeCalls.Load(); got != 0 {
		t.Fatalf("Decode calls = %d, want 0", got)
	}
}
