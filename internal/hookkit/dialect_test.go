package hookkit

import (
	"context"
	"errors"
	"testing"
)

type registerTestResults struct{}

func testDialect(t *testing.T) *Dialect {
	t.Helper()
	return NewDialect(NewCodec("test", errors.New("empty"), errors.New("decode"), errors.New("name")))
}

func TestRegisterWith_invokesWithResults(t *testing.T) {
	d := testDialect(t)

	var sawResults registerTestResults
	called := false
	RegisterWith(d, registerTestResults{}, func(_ context.Context, _ Hook[handlerTestEvent], r registerTestResults) (handlerTestOutput, error) {
		called = true
		sawResults = r
		return handlerTestOutput{body: "ok"}, nil
	})

	handlers := d.HandlersFor("HandlerTest")
	if len(handlers) != 1 {
		t.Fatalf("handlers = %d", len(handlers))
	}
	out, err := handlers[0].Invoke(context.Background(), handlerTestEvent{})
	if err != nil {
		t.Fatalf("err = %v", err)
	}
	got, ok := out.(handlerTestOutput)
	if !ok || got.body != "ok" || !called || sawResults != (registerTestResults{}) {
		t.Fatalf("got out=%#v called=%v results=%#v", out, called, sawResults)
	}
}

func TestRegisterWith_nilFn(t *testing.T) {
	d := testDialect(t)
	RegisterWith[handlerTestEvent, handlerTestOutput, registerTestResults](d, registerTestResults{}, nil)
	if got := d.HandlersFor("HandlerTest"); len(got) != 0 {
		t.Fatalf("handlers = %d", len(got))
	}
}

func TestRegisterObserve_invokes(t *testing.T) {
	d := testDialect(t)

	called := false
	RegisterObserve(d, func(_ context.Context, _ Hook[handlerTestEvent]) error {
		called = true
		return nil
	})

	handlers := d.HandlersFor("HandlerTest")
	if len(handlers) != 1 {
		t.Fatalf("handlers = %d", len(handlers))
	}
	out, err := handlers[0].Invoke(context.Background(), handlerTestEvent{})
	if err != nil || out != nil || !called {
		t.Fatalf("got out=%#v err=%v called=%v", out, err, called)
	}
}

func TestRegisterObserve_nilFn(t *testing.T) {
	d := testDialect(t)
	RegisterObserve[handlerTestEvent](d, nil)
	if got := d.HandlersFor("HandlerTest"); len(got) != 0 {
		t.Fatalf("handlers = %d", len(got))
	}
}
