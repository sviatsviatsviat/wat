package hooktest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ExitExpectFailed is returned when an optional expect document does not match
// the hook run. Expect mode maps a matching run to ExitOK regardless of the
// hook binary's own exit code.
const ExitExpectFailed = 1

// Expect is the optional assertion document loaded alongside a fixture.
//
// When an expect document is present, wat test compares the listed fields and
// returns ExitOK on success or ExitExpectFailed on mismatch. Omitted fields are
// not checked.
type Expect struct {
	// Exit, when set, must equal the hook process exit code.
	Exit *int `json:"exit"`
	// Decision, when non-empty, must equal the recognized decision field from
	// hook stdout (for example "deny" or "allow").
	Decision string `json:"decision"`
	// StdoutContains lists substrings that must appear in hook stdout.
	StdoutContains []string `json:"stdout_contains"`
	// StdoutEmpty, when set, asserts whether hook stdout is empty after trim.
	StdoutEmpty *bool `json:"stdout_empty"`
}

// SidecarExpectPath returns the conventional expect path for a fixture file:
// foo.json -> foo.expect.json. It returns "" for stdin ("-") or empty paths.
func SidecarExpectPath(fixture string) string {
	fixture = strings.TrimSpace(fixture)
	if fixture == "" || fixture == "-" {
		return ""
	}
	ext := filepath.Ext(fixture)
	if strings.EqualFold(ext, ".json") {
		return strings.TrimSuffix(fixture, ext) + ".expect.json"
	}
	return fixture + ".expect.json"
}

// ResolveExpectPath chooses the expect document path. An explicit --expect path
// wins. Otherwise the conventional sidecar is used when it exists.
func ResolveExpectPath(fixture, explicit string, stat func(string) (os.FileInfo, error)) (string, error) {
	explicit = strings.TrimSpace(explicit)
	if explicit != "" {
		if _, err := stat(explicit); err != nil {
			return "", fmt.Errorf("read expect: %w", err)
		}
		return explicit, nil
	}
	sidecar := SidecarExpectPath(fixture)
	if sidecar == "" {
		return "", nil
	}
	if _, err := stat(sidecar); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("stat expect: %w", err)
	}
	return sidecar, nil
}

// LoadExpect decodes an Expect document from path. Unknown JSON fields are
// rejected so typos fail closed.
func LoadExpect(path string, readFile func(string) ([]byte, error)) (Expect, error) {
	data, err := readFile(path)
	if err != nil {
		return Expect{}, fmt.Errorf("read expect: %w", err)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var exp Expect
	if err := dec.Decode(&exp); err != nil {
		return Expect{}, fmt.Errorf("decode expect: %w", err)
	}
	if dec.More() {
		return Expect{}, fmt.Errorf("decode expect: trailing data after JSON value")
	}
	return exp, nil
}

// CheckExpect compares hook stdout and exit against exp. It returns one
// human-readable mismatch line per failed assertion.
func CheckExpect(exp Expect, dialectName string, hookStdout []byte, hookExit int) []string {
	var fails []string
	if exp.Exit != nil && hookExit != *exp.Exit {
		fails = append(fails, fmt.Sprintf("exit: got %d, want %d", hookExit, *exp.Exit))
	}
	if exp.Decision != "" {
		got := summarizeHookDecision(dialectName, hookStdout)
		if got != exp.Decision {
			fails = append(fails, fmt.Sprintf("decision: got %q, want %q", got, exp.Decision))
		}
	}
	stdout := string(hookStdout)
	for _, want := range exp.StdoutContains {
		if !strings.Contains(stdout, want) {
			fails = append(fails, fmt.Sprintf("stdout_contains: missing %q", want))
		}
	}
	if exp.StdoutEmpty != nil {
		empty := len(bytes.TrimSpace(hookStdout)) == 0
		if empty != *exp.StdoutEmpty {
			if *exp.StdoutEmpty {
				fails = append(fails, "stdout_empty: want empty stdout")
			} else {
				fails = append(fails, "stdout_empty: want non-empty stdout")
			}
		}
	}
	return fails
}
