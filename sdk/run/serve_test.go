package run

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

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

type otherOutput struct{ body string }

func (o otherOutput) IsZero() bool                 { return o.body == "" }
func (o otherOutput) Encode() ([]byte, int, error) { return []byte(o.body), 0, nil }
func (o otherOutput) Stop() bool                   { return false }
func (o otherOutput) Merge(other Output) (Output, []string, error) {
	return nil, nil, fmt.Errorf("merge type mismatch: want otherOutput, got %T", other)
}

func newTestRouter(name, eventName string, decodeCalls *atomic.Int32) (*hookkit.Router, *hookkit.Dialect) {
	c := hookkit.NewCodec(name, fmt.Errorf("empty"), fmt.Errorf("decode"), fmt.Errorf("name required"))
	c.Register(eventName, func(raw []byte) (hookkit.Event, error) {
		if decodeCalls != nil {
			decodeCalls.Add(1)
		}
		return testEvent{raw: string(raw)}, nil
	})
	// For EventName peek, hookkit.Codec requires hook_event_name in JSON.
	// Tests pass payloads with that field or we register a custom path.
	d := hookkit.NewDialect(c)
	r := hookkit.NewRouter()
	r.Ensure(name, func([]byte, func(string) string) bool { return true }, d)
	return r, d
}

func TestServe_DecodesOnce(t *testing.T) {
	var decodeCalls atomic.Int32
	r, d := newTestRouter("testdialect", "TestEvent", &decodeCalls)
	var handlerCalls atomic.Int32
	for range 3 {
		d.Register(hookkit.Handler(func(_ context.Context, hook Hook[testEvent]) (emptyOutput, error) {
			handlerCalls.Add(1)
			if hook.Event.raw != `{"hook_event_name":"TestEvent","ok":true}` {
				t.Errorf("event = %#v", hook.Event)
			}
			return emptyOutput{}, nil
		}))
	}

	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent","ok":true}`), &bytes.Buffer{}, &bytes.Buffer{}, applyOptions())
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
	var decodeCalls atomic.Int32
	r, _ := newTestRouter("empty", "NoHandlers", &decodeCalls)
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"NoHandlers"}`), &bytes.Buffer{}, &bytes.Buffer{}, applyOptions())
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if got := decodeCalls.Load(); got != 0 {
		t.Fatalf("Decode calls = %d, want 0", got)
	}
}

func TestServe_FoldsOutputsEncodeOnce(t *testing.T) {
	var encodes atomic.Int32
	r, d := newTestRouter("fold", "TestEvent", nil)
	d.Register(hookkit.Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		return foldOutput{body: "first", encodes: &encodes}, nil
	}))
	d.Register(hookkit.Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		return foldOutput{body: "second", encodes: &encodes, exitCode: 2}, nil
	}))

	var stdout, stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &stdout, &stderr, applyOptions())
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
	r, d := newTestRouter("stop", "TestEvent", nil)
	var calls atomic.Int32
	d.Register(hookkit.Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		calls.Add(1)
		return foldOutput{body: "deny", stop: true}, nil
	}))
	d.Register(hookkit.Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		calls.Add(1)
		return foldOutput{body: "later"}, nil
	}))

	var stdout bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &stdout, &bytes.Buffer{}, applyOptions())
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
	r, d := newTestRouter("mismatch", "TestEvent", nil)
	d.Register(hookkit.Handler(func(context.Context, Hook[testEvent]) (foldOutput, error) {
		return foldOutput{body: "a"}, nil
	}))
	d.Register(hookkit.Handler(func(context.Context, Hook[testEvent]) (otherOutput, error) {
		return otherOutput{body: "b"}, nil
	}))

	var stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &bytes.Buffer{}, &stderr, applyOptions())
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "merge type mismatch") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}
