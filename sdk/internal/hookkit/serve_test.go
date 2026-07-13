package hookkit

import (
	"context"
	"io"
	"strings"
	"testing"
)

func TestServeLoop(t *testing.T) {
	t.Parallel()
	var outBuf, errBuf strings.Builder
	code := ServeLoop(context.Background(), strings.NewReader("payload"), &outBuf, &errBuf, ServeHooks{
		Label:            "test",
		HandlerErrorExit: 1,
		Decode: func([]byte) (string, any, error) {
			return "ev", "event", nil
		},
		Lookup: func(string) (func(context.Context, any) (any, error), bool) {
			return func(context.Context, any) (any, error) {
				return map[string]string{"ok": "1"}, nil
			}, true
		},
		Encode: func(string, any) ([]byte, int, error) {
			return []byte(`{"ok":1}`), 0, nil
		},
	})
	if code != 0 {
		t.Fatalf("exit code = %d", code)
	}
	if outBuf.String() != `{"ok":1}` {
		t.Fatalf("stdout = %q", outBuf.String())
	}
}

func TestServeLoopUnregisteredEvent(t *testing.T) {
	t.Parallel()
	code := ServeLoop(context.Background(), strings.NewReader("payload"), io.Discard, io.Discard, ServeHooks{
		Label:            "test",
		HandlerErrorExit: 1,
		Decode: func([]byte) (string, any, error) {
			return "ev", "event", nil
		},
		Lookup: func(string) (func(context.Context, any) (any, error), bool) {
			return nil, false
		},
	})
	if code != 0 {
		t.Fatalf("exit code = %d, want 0 for unregistered event", code)
	}
}
