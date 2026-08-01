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

type forcedEvt struct {
	raw string
}

func (forcedEvt) EventName() string { return "Forced" }

type emptyOutput struct{}

func (emptyOutput) IsZero() bool { return true }

func (emptyOutput) Encode() ([]byte, int, error) { return nil, 0, nil }

func (emptyOutput) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	return emptyOutput{}, nil, nil
}

func (emptyOutput) Stop() bool { return false }

type foldOutput struct {
	body         string
	stop         bool
	encodes      *atomic.Int32
	exitCode     int
	bodyOnStderr bool
}

func (o foldOutput) IsZero() bool { return o.body == "" && !o.stop }

func (o foldOutput) Encode() ([]byte, int, error) {
	if o.encodes != nil {
		o.encodes.Add(1)
	}
	return []byte(o.body), o.exitCode, nil
}

func (o foldOutput) BodyOnStderr() bool { return o.bodyOnStderr }

func (o foldOutput) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
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
	out.bodyOnStderr = o.bodyOnStderr || b.bodyOnStderr
	return out, warnings, nil
}

func (o foldOutput) Stop() bool { return o.stop }

type otherOutput struct{ body string }

func (o otherOutput) IsZero() bool                 { return o.body == "" }
func (o otherOutput) Encode() ([]byte, int, error) { return []byte(o.body), 0, nil }
func (o otherOutput) Stop() bool                   { return false }
func (o otherOutput) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	return nil, nil, fmt.Errorf("merge type mismatch: want otherOutput, got %T", other)
}

func newTestRouter(name, eventName string, decodeCalls *atomic.Int32) (*router, *hookkit.Dialect) {
	c := hookkit.NewCodec(name, fmt.Errorf("empty"), fmt.Errorf("decode"), fmt.Errorf("name required"))
	c.Register(eventName, func(raw []byte) (hookkit.Event, error) {
		if decodeCalls != nil {
			decodeCalls.Add(1)
		}
		return testEvent{raw: string(raw)}, nil
	})
	// For EventName peek, hookkit.Codec requires hook_event_name in JSON.
	// Tests pass payloads with that field or we register a custom path.
	r := newRouter()
	d := r.Ensure(name, func([]byte) bool { return true }, c)
	return r, d
}

func TestServe_DecodesOnce(t *testing.T) {
	var decodeCalls atomic.Int32
	r, d := newTestRouter("testdialect", "TestEvent", &decodeCalls)
	var handlerCalls atomic.Int32
	for range 3 {
		d.Register(hookkit.Handler(func(_ context.Context, hook testEvent) (emptyOutput, error) {
			handlerCalls.Add(1)
			if hook.raw != `{"hook_event_name":"TestEvent","ok":true}` {
				t.Errorf("event = %#v", hook)
			}
			return emptyOutput{}, nil
		}))
	}

	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent","ok":true}`), &bytes.Buffer{}, &bytes.Buffer{}, serveHints{})
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
	var stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"NoHandlers"}`), &bytes.Buffer{}, &stderr, serveHints{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if got := decodeCalls.Load(); got != 0 {
		t.Fatalf("Decode calls = %d, want 0", got)
	}
	if !strings.Contains(stderr.String(), `warning: no handlers for "NoHandlers"`) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_FoldsOutputsEncodeOnce(t *testing.T) {
	var encodes atomic.Int32
	r, d := newTestRouter("fold", "TestEvent", nil)
	d.Register(hookkit.Handler(func(context.Context, testEvent) (foldOutput, error) {
		return foldOutput{body: "first", encodes: &encodes}, nil
	}))
	d.Register(hookkit.Handler(func(context.Context, testEvent) (foldOutput, error) {
		return foldOutput{body: "second", encodes: &encodes, exitCode: 2}, nil
	}))

	var stdout, stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &stdout, &stderr, serveHints{})
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

func TestServe_BodyOnStderrWritesStderr(t *testing.T) {
	r, d := newTestRouter("any", "TestEvent", nil)
	d.Register(hookkit.Handler(func(context.Context, testEvent) (foldOutput, error) {
		return foldOutput{body: "keep working", exitCode: 2, bodyOnStderr: true}, nil
	}))

	var stdout, stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &stdout, &stderr, serveHints{})
	if code != 2 {
		t.Fatalf("exit = %d, want 2", code)
	}
	if stdout.Len() != 0 {
		t.Fatalf("stdout = %q, want empty when BodyOnStderr", stdout.String())
	}
	if got := strings.TrimSpace(stderr.String()); got != "keep working" {
		t.Fatalf("stderr = %q, want keep working", stderr.String())
	}
}

func TestServe_StopSkipsLaterHandlers(t *testing.T) {
	r, d := newTestRouter("stop", "TestEvent", nil)
	var calls atomic.Int32
	d.Register(hookkit.Handler(func(context.Context, testEvent) (foldOutput, error) {
		calls.Add(1)
		return foldOutput{body: "deny", stop: true}, nil
	}))
	d.Register(hookkit.Handler(func(context.Context, testEvent) (foldOutput, error) {
		calls.Add(1)
		return foldOutput{body: "later"}, nil
	}))

	var stdout bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &stdout, &bytes.Buffer{}, serveHints{})
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
	d.Register(hookkit.Handler(func(context.Context, testEvent) (foldOutput, error) {
		return foldOutput{body: "a"}, nil
	}))
	d.Register(hookkit.Handler(func(context.Context, testEvent) (otherOutput, error) {
		return otherOutput{body: "b"}, nil
	}))

	var stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &bytes.Buffer{}, &stderr, serveHints{})
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), "merge type mismatch") {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_EventHintSkipsPeek(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	r := newRouter()
	c := hookkit.NewCodec("hint", fmt.Errorf("empty"), fmt.Errorf("decode"), fmt.Errorf("name required"))
	c.Register("Forced", func(raw []byte) (hookkit.Event, error) {
		return forcedEvt{raw: string(raw)}, nil
	})
	d := r.Ensure("hint", func([]byte) bool { return true }, c)
	d.Register(hookkit.Handler(func(_ context.Context, hook forcedEvt) (emptyOutput, error) {
		calls.Add(1)
		if !strings.Contains(hook.raw, `"other":true`) {
			t.Errorf("raw = %q", hook.raw)
		}
		return emptyOutput{}, nil
	}))

	var stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"other":true}`), &bytes.Buffer{}, &stderr, serveHints{event: "Forced"})
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if calls.Load() != 1 {
		t.Fatalf("handler calls = %d, want 1", calls.Load())
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestServe_EventHintMismatchWarns(t *testing.T) {
	t.Parallel()
	var calls atomic.Int32
	r := newRouter()
	c := hookkit.NewCodec("hint", fmt.Errorf("empty"), fmt.Errorf("decode"), fmt.Errorf("name required"))
	c.Register("Forced", func(raw []byte) (hookkit.Event, error) {
		return forcedEvt{raw: string(raw)}, nil
	})
	d := r.Ensure("hint", func([]byte) bool { return true }, c)
	d.Register(hookkit.Handler(func(context.Context, forcedEvt) (emptyOutput, error) {
		calls.Add(1)
		return emptyOutput{}, nil
	}))

	var stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"Other"}`), &bytes.Buffer{}, &stderr, serveHints{event: "Forced"})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if calls.Load() != 1 {
		t.Fatalf("handler calls = %d, want 1", calls.Load())
	}
	if !strings.Contains(stderr.String(), `--event "Forced" disagrees with hook_event_name "Other"`) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_AgentHintLookup(t *testing.T) {
	t.Parallel()
	r := newRouter()
	c := hookkit.NewCodec("forced", fmt.Errorf("empty"), fmt.Errorf("decode"), fmt.Errorf("name required"))
	c.Register("TestEvent", func(raw []byte) (hookkit.Event, error) {
		return testEvent{raw: string(raw)}, nil
	})
	var detectHits atomic.Int32
	d := r.Ensure("forced", func([]byte) bool {
		detectHits.Add(1)
		return false
	}, c)
	var calls atomic.Int32
	d.Register(hookkit.Handler(func(context.Context, testEvent) (emptyOutput, error) {
		calls.Add(1)
		return emptyOutput{}, nil
	}))
	// A second dialect that would match detection if hint were absent.
	other := hookkit.NewCodec("other", fmt.Errorf("empty"), fmt.Errorf("decode"), fmt.Errorf("name required"))
	other.Register("TestEvent", func(raw []byte) (hookkit.Event, error) {
		return testEvent{raw: string(raw)}, nil
	})
	r.Ensure("other", func([]byte) bool { return true }, other)

	var stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &bytes.Buffer{}, &stderr, serveHints{agent: "forced"})
	if code != 0 {
		t.Fatalf("exit = %d, stderr = %q", code, stderr.String())
	}
	if calls.Load() != 1 {
		t.Fatalf("handler calls = %d, want 1", calls.Load())
	}
	if !strings.Contains(stderr.String(), `--agent "forced" disagrees with detected dialect "other"`) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestServe_UnknownAgentHint(t *testing.T) {
	t.Parallel()
	r, _ := newTestRouter("only", "TestEvent", nil)
	var stderr bytes.Buffer
	code := serve(context.Background(), r, strings.NewReader(`{"hook_event_name":"TestEvent"}`), &bytes.Buffer{}, &stderr, serveHints{agent: "missing"})
	if code != 1 {
		t.Fatalf("exit = %d, want 1", code)
	}
	if !strings.Contains(stderr.String(), `unknown dialect "missing"`) {
		t.Fatalf("stderr = %q", stderr.String())
	}
}

func TestParseServeHints(t *testing.T) {
	t.Parallel()
	got := parseServeHints([]string{"--agent", "claude", "--event", "PreToolUse"})
	if got.agent != "claude" || got.event != "PreToolUse" {
		t.Fatalf("got = %#v", got)
	}
	got = parseServeHints([]string{"--event", "SessionStart"})
	if got.agent != "" || got.event != "SessionStart" {
		t.Fatalf("got = %#v", got)
	}
	got = parseServeHints(nil)
	if got != (serveHints{}) {
		t.Fatalf("got = %#v", got)
	}
}
