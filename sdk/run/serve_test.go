package run

import (
	"bytes"
	"context"
	"fmt"
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

func (emptyOutput) Merge(other Output) (Output, []string, error) {
	return emptyOutput{}, nil, nil
}

func (emptyOutput) Stop() bool { return false }

type foldOutput struct {
	body     string
	stop     bool
	encodes  *atomic.Int32
	exitCode int
}

func (o foldOutput) IsZero() bool { return o.body == "" && !o.stop }

func (o foldOutput) Encode() ([]byte, int, error) {
	if o.encodes != nil {
		o.encodes.Add(1)
	}
	return []byte(o.body), o.exitCode, nil
}

func (o foldOutput) Merge(other Output) (Output, []string, error) {
	b, ok := other.(foldOutput)
	if !ok {
		return nil, nil, fmt.Errorf("merge type mismatch: want foldOutput, got %T", other)
	}
	out := o
	var warnings []string
	if o.body != "" && b.body != "" {
		warnings = append(warnings, "body: overwritten by later handler")
	}
	if b.body != "" {
		out.body = b.body
	}
	if b.encodes != nil {
		out.encodes = b.encodes
	}
	if b.exitCode > out.exitCode {
		out.exitCode = b.exitCode
	}
	out.stop = o.stop || b.stop
	return out, warnings, nil
}

func (o foldOutput) Stop() bool { return o.stop }

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

func TestServe_FoldsOutputsEncodeOnce(t *testing.T) {
	r := NewRegistry()
	var encodes atomic.Int32
	r.registerDialect("fold", DialectOps{
		Detect: func([]byte, func(string) string) bool { return true },
		Codec:  testCodec{eventName: "TestEvent"},
	})
	r.RegisterHandler("fold", Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		return foldOutput{body: "first", encodes: &encodes}, nil
	}))
	r.RegisterHandler("fold", Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		return foldOutput{body: "second", encodes: &encodes, exitCode: 2}, nil
	}))

	var stdout, stderr bytes.Buffer
	code := r.serve(context.Background(), strings.NewReader(`{}`), &stdout, &stderr, applyOptions())
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if stdout.String() != "second" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if encodes.Load() != 1 {
		t.Fatalf("Encode calls = %d, want 1", encodes.Load())
	}
	if !strings.Contains(stderr.String(), "body: overwritten by later handler") {
		t.Fatalf("stderr = %q, want overwrite warning", stderr.String())
	}
}

func TestServe_StopSkipsLaterHandlers(t *testing.T) {
	r := NewRegistry()
	r.registerDialect("stop", DialectOps{
		Detect: func([]byte, func(string) string) bool { return true },
		Codec:  testCodec{eventName: "TestEvent"},
	})
	var calls atomic.Int32
	r.RegisterHandler("stop", Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		calls.Add(1)
		return foldOutput{body: "deny", stop: true}, nil
	}))
	r.RegisterHandler("stop", Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		calls.Add(1)
		return foldOutput{body: "later"}, nil
	}))

	var stdout bytes.Buffer
	code := r.serve(context.Background(), strings.NewReader(`{}`), &stdout, &bytes.Buffer{}, applyOptions())
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if stdout.String() != "deny" {
		t.Fatalf("stdout = %q", stdout.String())
	}
	if calls.Load() != 1 {
		t.Fatalf("handler calls = %d, want 1", calls.Load())
	}
}

func TestServe_MergeTypeMismatch(t *testing.T) {
	r := NewRegistry()
	r.registerDialect("mismatch", DialectOps{
		Detect: func([]byte, func(string) string) bool { return true },
		Codec:  testCodec{eventName: "TestEvent"},
	})
	r.RegisterHandler("mismatch", Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		return foldOutput{body: "a"}, nil
	}))
	r.RegisterHandler("mismatch", Handler(func(context.Context, Hook[testEvent]) (handlerTestOutput, error) {
		return handlerTestOutput{body: "b"}, nil
	}))

	var stderr bytes.Buffer
	code := r.serve(context.Background(), strings.NewReader(`{}`), &bytes.Buffer{}, &stderr, applyOptions())
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "merge type mismatch") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
