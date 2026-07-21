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
}

func (o handlerTestOutput) IsZero() bool { return o.body == "" }

func (o handlerTestOutput) Encode() ([]byte, int, error) {
	return []byte(o.body), 0, nil
}

func TestHandler_handle(t *testing.T) {
	reg := Handler(func(_ context.Context, hook Hook[handlerTestEvent]) (handlerTestOutput, error) {
		return handlerTestOutput{body: "ok"}, nil
	})
	if reg.eventName() != "HandlerTest" {
		t.Fatalf("eventName = %q", reg.eventName())
	}
	out, code, err := reg.handle(context.Background(), handlerTestEvent{})
	if err != nil || code != 0 || string(out) != "ok" {
		t.Fatalf("got %q code=%d err=%v", out, code, err)
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
	_, _, _ = reg.handle(context.Background(), testEvent{raw: "other"})
	t.Fatal("want panic")
}

func TestHandler_encodeOutput_zero(t *testing.T) {
	out, code, err := encodeOutput(handlerTestOutput{})
	if err != nil || code != 0 || out != nil {
		t.Fatalf("got %q code=%d err=%v", out, code, err)
	}
}

func TestObserveHandler_handle(t *testing.T) {
	called := false
	reg := ObserveHandler(func(_ context.Context, hook Hook[handlerTestEvent]) error {
		called = true
		return nil
	})
	out, code, err := reg.handle(context.Background(), handlerTestEvent{})
	if err != nil || code != 0 || out != nil || !called {
		t.Fatalf("got out=%q code=%d err=%v called=%v", out, code, err, called)
	}
}

func TestObserveHandler_error(t *testing.T) {
	want := errors.New("boom")
	reg := ObserveHandler(func(context.Context, Hook[handlerTestEvent]) error {
		return want
	})
	_, code, err := reg.handle(context.Background(), handlerTestEvent{})
	if code != 1 || !errors.Is(err, want) {
		t.Fatalf("code=%d err=%v", code, err)
	}
}
