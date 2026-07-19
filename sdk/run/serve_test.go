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

func (c testCodec) EventName([]byte) (string, error) {
	return c.eventName, nil
}

func (c testCodec) Decode(raw []byte) (any, error) {
	if c.decodeCalls != nil {
		c.decodeCalls.Add(1)
	}
	return string(raw), nil
}

func TestServe_DecodesOnce(t *testing.T) {
	r := NewRegistry()
	var decodeCalls atomic.Int32
	r.RegisterDialect("testdialect", DialectOps{
		Detect: func([]byte, func(string) string) bool { return true },
		Codec:  testCodec{eventName: "TestEvent", decodeCalls: &decodeCalls},
	})
	var handlerCalls atomic.Int32
	for range 3 {
		r.RegisterHandler("testdialect", "TestEvent", func(_ context.Context, event any) ([]byte, int, error) {
			handlerCalls.Add(1)
			if event != `{"ok":true}` {
				t.Errorf("event = %#v", event)
			}
			return nil, 0, nil
		})
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
	r.RegisterDialect("empty", DialectOps{
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
