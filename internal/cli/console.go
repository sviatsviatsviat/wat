// Package cli provides console I/O for diagnostics and hook protocol output, help text, exit codes, and hook stdin parsing.
package cli

import (
	"fmt"
	"io"
)

// Console separates diagnostic output (errors, help) from hook protocol output on stdout.
type Console interface {
	// WriteError writes args to the diagnostic stream as one println-style line: arguments are
	// space-separated and followed by a newline (same as fmt.Fprintln).
	WriteError(args ...any) error
	// WriteErrorf writes a formatted message to the diagnostic stream using fmt.Fprintf. It does
	// not append a newline; include "\n" in format (or in a final %s argument) when you want line
	// breaks. Most call sites use a format string ending with "\n" for a single diagnostic line.
	WriteErrorf(format string, args ...any) error
	// Write writes s to the hook protocol stream without adding a newline.
	Write(s string) error
}

// dualStreamConsole implements [Console] with separate diagnostic and hook writers.
type dualStreamConsole struct {
	stderr     io.Writer
	hookStdout io.Writer
}

func (c dualStreamConsole) writeDiagnosticLine(args ...any) error {
	_, err := fmt.Fprintln(c.stderr, args...)
	return err
}

func (c dualStreamConsole) writeDiagnosticf(format string, args ...any) error {
	_, err := fmt.Fprintf(c.stderr, format, args...)
	return err
}

func (c dualStreamConsole) WriteError(args ...any) error {
	return c.writeDiagnosticLine(args...)
}

func (c dualStreamConsole) WriteErrorf(format string, args ...any) error {
	return c.writeDiagnosticf(format, args...)
}

func (c dualStreamConsole) Write(s string) error {
	_, err := io.WriteString(c.hookStdout, s)
	return err
}

// NewConsole returns a [Console] that writes diagnostics to stderr and hook output to hookStdout.
// Nil writers are replaced with [io.Discard] so [Console] methods never panic on nil [io.Writer].
func NewConsole(stderr, hookStdout io.Writer) Console {
	if stderr == nil {
		stderr = io.Discard
	}
	if hookStdout == nil {
		hookStdout = io.Discard
	}
	return dualStreamConsole{stderr: stderr, hookStdout: hookStdout}
}
