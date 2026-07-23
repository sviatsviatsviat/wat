package runtime

import (
	"fmt"
	"os"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
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
		if !hookkit.ValidEnvKey(k) {
			return fmt.Errorf("claude: invalid env key %q", k)
		}
		buf = fmt.Appendf(buf, "export %s=%s\n", k, hookkit.ShellSingleQuote(v))
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
