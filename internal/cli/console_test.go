package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestNewConsole_WriteError(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf, io.Discard)
	if err := console.WriteError("one", 2); err != nil {
		t.Fatalf("WriteError: %v", err)
	}
	if written := buf.String(); written != "one 2\n" {
		t.Fatalf("WriteError output: got %q want %q", written, "one 2\n")
	}
}

func TestNewConsole_WriteErrorf(t *testing.T) {
	var buf bytes.Buffer
	console := NewConsole(&buf, io.Discard)
	if err := console.WriteErrorf("x=%d\n", 7); err != nil {
		t.Fatalf("WriteErrorf: %v", err)
	}
	if written := buf.String(); written != "x=7\n" {
		t.Fatalf("WriteErrorf output: got %q want %q", written, "x=7\n")
	}
}

func TestNewConsole_nilWritersDoNotPanic(t *testing.T) {
	c := NewConsole(nil, nil)
	if err := c.WriteError("x"); err != nil {
		t.Fatalf("WriteError: %v", err)
	}
	if err := c.WriteErrorf("y\n"); err != nil {
		t.Fatalf("WriteErrorf: %v", err)
	}
	if err := c.Write("z"); err != nil {
		t.Fatalf("Write: %v", err)
	}
}

func TestNewConsole_Write_stdout(t *testing.T) {
	var stderr, stdout bytes.Buffer
	console := NewConsole(&stderr, &stdout)
	if err := console.Write("{}\n"); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if want := "{}\n"; stdout.String() != want {
		t.Fatalf("stdout: got %q want %q", stdout.String(), want)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr should be empty, got %q", stderr.String())
	}
}

func TestNewConsole_propagatesWriteError(t *testing.T) {
	console := NewConsole(errWriter{}, io.Discard)
	err := console.WriteError("fail")
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected write error, got %v", err)
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) {
	return 0, errWriteFailed{}
}

type errWriteFailed struct{}

func (errWriteFailed) Error() string { return "boom" }
