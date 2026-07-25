package portio

import (
	"bytes"
	"errors"
	"io"
	"testing"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
)

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestRun_stdoutWriteError(t *testing.T) {
	deps := DefaultDeps()
	deps.ReadFile = func(string) ([]byte, error) {
		return []byte(`{"hooks":{}}`), nil
	}
	var errBuf bytes.Buffer
	code := Run(Config{
		From:      sdkclaude.Dialect,
		To:        sdkclaude.Dialect,
		InputPath: "in.json",
	}, deps, errWriter{}, &errBuf)
	if code != ExitRuntimeFailure {
		t.Fatalf("exit = %d, want %d", code, ExitRuntimeFailure)
	}
	if !bytes.Contains(errBuf.Bytes(), []byte("write stdout")) {
		t.Fatalf("stderr = %q", errBuf.String())
	}
}

func TestRun_stdoutOK(t *testing.T) {
	deps := DefaultDeps()
	in := []byte(`{"hooks":{}}`)
	deps.ReadFile = func(string) ([]byte, error) { return in, nil }
	var out, errBuf bytes.Buffer
	code := Run(Config{
		From:      sdkclaude.Dialect,
		To:        sdkclaude.Dialect,
		InputPath: "in.json",
	}, deps, &out, &errBuf)
	if code != ExitOK {
		t.Fatalf("exit = %d, want %d: %s", code, ExitOK, errBuf.String())
	}
	if !bytes.HasSuffix(out.Bytes(), []byte("\n")) {
		t.Fatalf("stdout should end with newline: %q", out.String())
	}
	_ = io.Discard
}
