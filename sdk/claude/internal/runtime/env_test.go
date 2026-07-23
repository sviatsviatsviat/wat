package runtime

import "testing"

func TestWriteEnvFile_InvalidKey(t *testing.T) {
	err := WriteEnvFile(
		map[string]string{"FOO\nBAR": "value"},
		func(string) string { return "/tmp/env.sh" },
		nil,
	)
	if err == nil {
		t.Fatal("expected error for invalid env key")
	}
}
