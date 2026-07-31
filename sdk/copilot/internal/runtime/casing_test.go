package runtime_test

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

func TestAdviseMissingHandlers_CamelCase(t *testing.T) {
	t.Parallel()
	err := runtime.AdviseMissingHandlers("preToolUse")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `PascalCase event "PreToolUse"`) {
		t.Fatalf("err = %v", err)
	}
}

func TestAdviseMissingHandlers_CaseFold(t *testing.T) {
	t.Parallel()
	err := runtime.AdviseMissingHandlers("PRETOOLUSE")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), `PascalCase event "PreToolUse"`) {
		t.Fatalf("err = %v", err)
	}
}

func TestAdviseMissingHandlers_PascalCaseUnhandled(t *testing.T) {
	t.Parallel()
	if err := runtime.AdviseMissingHandlers("SessionStart"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := runtime.AdviseMissingHandlers("UnknownEvent"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterCasingRejects_DecodeAs(t *testing.T) {
	t.Parallel()
	c := hookkit.NewCodec(runtime.Dialect, runtime.ErrEmptyPayload, runtime.ErrDecodePayload, runtime.ErrEventNameRequired)
	runtime.RegisterCasingRejects(c)
	c.Register("PreToolUse", func([]byte) (hookkit.Event, error) {
		t.Fatal("PascalCase decoder must not run for camelCase name")
		return nil, nil
	})
	_, err := c.DecodeAs([]byte(`{}`), "preToolUse")
	if err == nil {
		t.Fatal("expected decode error")
	}
	if !strings.Contains(err.Error(), `use PascalCase "PreToolUse"`) {
		t.Fatalf("err = %v", err)
	}
}
