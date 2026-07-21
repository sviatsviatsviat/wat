package run

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type handlerTestEvent struct{}

func (handlerTestEvent) EventName() string { return "HandlerTest" }

type handlerTestOutput struct {
	body string
	stop bool
}

func (o handlerTestOutput) IsZero() bool { return o.body == "" && !o.stop }

func (o handlerTestOutput) Encode() ([]byte, int, error) {
	return []byte(o.body), 0, nil
}

func (o handlerTestOutput) Merge(other Output) (Output, []string, error) {
	b, ok := other.(handlerTestOutput)
	if !ok {
		return nil, nil, fmt.Errorf("merge type mismatch: want handlerTestOutput, got %T", other)
	}
	out := o
	var warnings []string
	if o.body != "" && b.body != "" {
		warnings = append(warnings, "body: overwritten by later handler")
	}
	if b.body != "" {
		out.body = b.body
	}
	out.stop = o.stop || b.stop
	return out, warnings, nil
}

func (o handlerTestOutput) Stop() bool { return o.stop }

func TestHandler_invoke(t *testing.T) {
	reg := Handler(func(_ context.Context, hook Hook[handlerTestEvent]) (handlerTestOutput, error) {
		return handlerTestOutput{body: "ok"}, nil
	})
	if reg.eventName() != "HandlerTest" {
		t.Fatalf("eventName = %q", reg.eventName())
	}
	out, err := reg.invoke(context.Background(), handlerTestEvent{})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	got, ok := out.(handlerTestOutput)
	if !ok || got.body != "ok" {
		t.Fatalf("got %#v", out)
	}
}

func TestHandler_typeMismatch(t *testing.T) {
	reg := Handler(func(context.Context, Hook[handlerTestEvent]) (handlerTestOutput, error) {
		return handlerTestOutput{}, nil
	})
	defer func() {
		got := fmt.Sprint(recover())
		if !strings.Contains(got, "handler for HandlerTest received") {
			t.Fatalf("panic = %q", got)
		}
	}()
	_, _ = reg.invoke(context.Background(), testEvent{raw: "other"})
	t.Fatal("want panic")
}

func TestHandler_zeroOutput(t *testing.T) {
	reg := Handler(func(context.Context, Hook[handlerTestEvent]) (handlerTestOutput, error) {
		return handlerTestOutput{}, nil
	})
	out, err := reg.invoke(context.Background(), handlerTestEvent{})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	if out == nil || !out.IsZero() {
		t.Fatalf("got %#v", out)
	}
}

func TestObserveHandler_invoke(t *testing.T) {
	called := false
	reg := ObserveHandler(func(_ context.Context, hook Hook[handlerTestEvent]) error {
		called = true
		return nil
	})
	out, err := reg.invoke(context.Background(), handlerTestEvent{})
	if err != nil || out != nil || !called {
		t.Fatalf("got out=%#v err=%v called=%v", out, err, called)
	}
}

func TestObserveHandler_error(t *testing.T) {
	want := errors.New("boom")
	reg := ObserveHandler(func(context.Context, Hook[handlerTestEvent]) error {
		return want
	})
	_, err := reg.invoke(context.Background(), handlerTestEvent{})
	if !errors.Is(err, want) {
		t.Fatalf("err=%v", err)
	}
}
