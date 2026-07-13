package claudehook

import (
	"fmt"
	"os"
	"strings"
)

// WriteEnvFile appends export lines for env to CLAUDE_ENV_FILE when set.
func WriteEnvFile(env map[string]string, getenv func(string) string, appendFile func(path string, data []byte) error) error {
	if len(env) == 0 {
		return nil
	}
	if getenv == nil {
		getenv = os.Getenv
	}
	path := getenv("CLAUDE_ENV_FILE")
	if path == "" {
		return nil
	}
	if appendFile == nil {
		appendFile = defaultAppendFile
	}
	var buf []byte
	for k, v := range env {
		if !validEnvKey(k) {
			return fmt.Errorf("claudehook: invalid env key %q", k)
		}
		buf = append(buf, []byte(fmt.Sprintf("export %s=%s\n", k, shellSingleQuote(v)))...)
	}
	return appendFile(path, buf)
}

func defaultAppendFile(path string, data []byte) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600) //nolint:gosec // CLAUDE_ENV_FILE path from agent env
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	_, err = f.Write(data)
	return err
}

// shellSingleQuote wraps s in single quotes using POSIX-safe escaping.
func shellSingleQuote(s string) string {
	if s == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}

// validEnvKey reports whether k is a valid ASCII shell export name.
func validEnvKey(k string) bool {
	if k == "" {
		return false
	}
	for i := 0; i < len(k); i++ {
		c := k[i]
		if i == 0 {
			if c != '_' && (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') {
				return false
			}
			continue
		}
		if c != '_' && (c < 'A' || c > 'Z') && (c < 'a' || c > 'z') && (c < '0' || c > '9') {
			return false
		}
	}
	return true
}
