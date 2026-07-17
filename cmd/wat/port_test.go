package main

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

const claudePortFixture = `{
  "hooks": {
    "PreToolUse": [
      {
        "matcher": "Bash",
        "hooks": [
          {"type": "command", "command": ".claude/hooks/block-rm.sh", "timeout": 15}
        ]
      },
      {
        "matcher": "Edit|Write",
        "hooks": [
          {"type": "command", "command": ".claude/hooks/lint.sh"}
        ]
      }
    ],
    "Stop": [
      {"hooks": [{"type": "command", "command": ".claude/hooks/require-tests.sh"}]}
    ],
    "MessageDisplay": [
      {"hooks": [{"type": "command", "command": ".claude/hooks/plain.sh"}]}
    ]
  }
}`

func TestPortProject_claudeToCursor(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "settings.json")
	if err := os.WriteFile(inputPath, []byte(claudePortFixture), 0o644); err != nil {
		t.Fatal(err)
	}

	prevStdout, prevStderr := stdout, stderr
	var outBuf, errBuf bytes.Buffer
	stdout, stderr = &outBuf, &errBuf
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	code := portProject(portConfig{
		from:      sdkclaude.Dialect,
		to:        sdkcursor.Dialect,
		inputPath: inputPath,
	}, defaultPortDeps())
	if code != exitOK {
		t.Fatalf("exit = %d, want %d\nstderr: %s", code, exitOK, errBuf.String())
	}

	var f struct {
		Version int `json:"version"`
		Hooks   map[string][]struct {
			Command string `json:"command"`
			Matcher string `json:"matcher"`
		} `json:"hooks"`
	}
	if err := json.Unmarshal(outBuf.Bytes(), &f); err != nil {
		t.Fatalf("parse stdout JSON: %v\n%s", err, outBuf.String())
	}
	if f.Version != 1 {
		t.Fatalf("version = %d, want 1", f.Version)
	}
	pre := f.Hooks["preToolUse"]
	if len(pre) != 2 {
		t.Fatalf("preToolUse handlers = %d, want 2: %s", len(pre), outBuf.String())
	}
	var matchers []string
	for _, h := range pre {
		matchers = append(matchers, h.Matcher)
	}
	if !strings.Contains(strings.Join(matchers, " "), "Shell") {
		t.Errorf("Bash should map to Shell, got %q", matchers)
	}
	if len(f.Hooks["stop"]) != 1 {
		t.Errorf("stop should have 1 handler: %s", outBuf.String())
	}
	if _, ok := f.Hooks["MessageDisplay"]; ok {
		t.Error("MessageDisplay must not appear in Cursor config")
	}
	if !strings.Contains(errBuf.String(), "MessageDisplay") {
		t.Errorf("stderr missing MessageDisplay warning: %q", errBuf.String())
	}
}

func TestPortProject_writesOutputFile(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "in.json")
	outputPath := filepath.Join(dir, "out.json")
	if err := os.WriteFile(inputPath, []byte(claudePortFixture), 0o644); err != nil {
		t.Fatal(err)
	}

	prevStdout, prevStderr := stdout, stderr
	stdout, stderr = io.Discard, io.Discard
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	code := portProject(portConfig{
		from:       sdkclaude.Dialect,
		to:         sdkcursor.Dialect,
		inputPath:  inputPath,
		outputPath: outputPath,
	}, defaultPortDeps())
	if code != exitOK {
		t.Fatalf("exit = %d, want %d", code, exitOK)
	}

	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if data[len(data)-1] != '\n' {
		t.Fatal("output file should end with newline")
	}
	var f struct {
		Version int `json:"version"`
	}
	if err := json.Unmarshal(data, &f); err != nil {
		t.Fatalf("parse output JSON: %v\n%s", err, data)
	}
	if f.Version != 1 {
		t.Fatalf("version = %d, want 1", f.Version)
	}
}

func TestPortProject_missingInput(t *testing.T) {
	prevStdout, prevStderr := stdout, stderr
	stdout, stderr = io.Discard, io.Discard
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	deps := defaultPortDeps()
	deps.readFile = func(string) ([]byte, error) {
		return nil, os.ErrNotExist
	}

	code := portProject(portConfig{
		from:      sdkclaude.Dialect,
		to:        sdkcursor.Dialect,
		inputPath: "/nonexistent/settings.json",
	}, deps)
	if code != exitRuntimeFailure {
		t.Fatalf("exit = %d, want %d", code, exitRuntimeFailure)
	}
}

func TestPortProject_sameDialect(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "settings.json")
	if err := os.WriteFile(inputPath, []byte(claudePortFixture), 0o644); err != nil {
		t.Fatal(err)
	}

	prevStdout, prevStderr := stdout, stderr
	var outBuf, errBuf bytes.Buffer
	stdout, stderr = &outBuf, &errBuf
	t.Cleanup(func() { stdout, stderr = prevStdout, prevStderr })

	code := portProject(portConfig{
		from:      sdkclaude.Dialect,
		to:        sdkclaude.Dialect,
		inputPath: inputPath,
	}, defaultPortDeps())
	if code != exitOK {
		t.Fatalf("exit = %d, want %d", code, exitOK)
	}
	if outBuf.String() != claudePortFixture+"\n" {
		t.Fatalf("same-dialect output mismatch:\n%s", outBuf.String())
	}
	if errBuf.Len() != 0 {
		t.Fatalf("unexpected warnings: %q", errBuf.String())
	}
}

func TestParsePortDialect(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		flag    string
		wantErr string
	}{
		{name: "empty", value: "", flag: "from", wantErr: "--from is required"},
		{name: "unknown", value: "nosuch", flag: "to", wantErr: "unknown agent dialect"},
		{name: "claude", value: "claude", flag: "from"},
		{name: "copilot", value: "copilot", flag: "to"},
		{name: "cursor", value: "cursor", flag: "from"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePortDialect(tt.value, tt.flag)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %q, want substring %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got == "" {
				t.Fatal("got Unknown dialect")
			}
		})
	}
}

func TestDefaultInputPath(t *testing.T) {
	wd := "/project"
	tests := []struct {
		from string
		want string
	}{
		{from: sdkclaude.Dialect, want: "/project/.claude/settings.json"},
		{from: sdkcopilot.Dialect, want: "/project/.github/hooks/wat.json"},
		{from: sdkcursor.Dialect, want: "/project/.cursor/hooks.json"},
	}
	for _, tt := range tests {
		t.Run(tt.from, func(t *testing.T) {
			if got := defaultInputPath(tt.from, wd); got != tt.want {
				t.Fatalf("defaultInputPath = %q, want %q", got, tt.want)
			}
		})
	}
}
