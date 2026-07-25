package e2e_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWatInstall_discoversAuthoredHooks(t *testing.T) {
	binary := buildWat(t)
	project := initProjectWithReplace(t)

	stdout, stderr, code := runWat(t, binary, project, "install", "--agent", "cursor", "--wat-path", binary)
	if code != 0 {
		t.Fatalf("exit = %d, want 0\nstdout:\n%s\nstderr:\n%s", code, stdout, stderr)
	}

	cursorPath := filepath.Join(project, ".cursor", "hooks.json")
	body, err := os.ReadFile(cursorPath)
	if err != nil {
		t.Fatal(err)
	}
	var installed struct {
		Version int                          `json:"version"`
		Hooks   map[string][]json.RawMessage `json:"hooks"`
	}
	if err := json.Unmarshal(body, &installed); err != nil {
		t.Fatal(err)
	}
	for _, event := range []string{
		"beforeShellExecution",
		"afterFileEdit",
		"stop",
	} {
		if len(installed.Hooks[event]) != 1 {
			t.Fatalf("%s handlers = %d, want 1", event, len(installed.Hooks[event]))
		}
	}
	if _, ok := installed.Hooks["sessionStart"]; ok {
		t.Fatalf("unregistered sessionStart event was installed")
	}
}
