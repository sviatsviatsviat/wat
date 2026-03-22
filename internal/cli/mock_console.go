package cli

import (
	"bytes"
	"io"
)

// MockConsole is a [Console] implementation that records diagnostic and hook streams in memory.
type MockConsole struct {
	stderrBuf bytes.Buffer
	stdoutBuf bytes.Buffer
	dualStreamConsole
}

// NewMockConsole returns a mock with empty buffers; [MockConsole.StderrBufferWriter] shares the diagnostic buffer.
func NewMockConsole() *MockConsole {
	mock := &MockConsole{}
	mock.dualStreamConsole = dualStreamConsole{stderr: &mock.stderrBuf, hookStdout: &mock.stdoutBuf}
	return mock
}

// StderrBufferWriter returns a writer that appends to the same buffer as [Console.WriteError] and [Console.WriteErrorf].
func (mock *MockConsole) StderrBufferWriter() io.Writer {
	return &mock.stderrBuf
}

// StdoutBufferWriter returns a writer that appends to the same buffer as [Console.Write].
func (mock *MockConsole) StdoutBufferWriter() io.Writer {
	return &mock.stdoutBuf
}

// StderrString returns recorded diagnostic output.
func (mock *MockConsole) StderrString() string {
	return mock.stderrBuf.String()
}

// StdoutString returns recorded hook protocol output.
func (mock *MockConsole) StdoutString() string {
	return mock.stdoutBuf.String()
}

// StdoutContains reports whether recorded hook output contains sub.
func (mock *MockConsole) StdoutContains(sub string) bool {
	return bytes.Contains(mock.stdoutBuf.Bytes(), []byte(sub))
}

// StderrContains reports whether recorded diagnostic output contains sub.
func (mock *MockConsole) StderrContains(sub string) bool {
	return bytes.Contains(mock.stderrBuf.Bytes(), []byte(sub))
}
