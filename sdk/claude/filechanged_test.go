package claude

import (
	"testing"
)

func TestDecode_FileChanged(t *testing.T) {
	ev := mustDecode[FileChanged](t, `{"session_id":"s","hook_event_name":"FileChanged","file_path":"/f.go"}`, EventFileChanged)
	if ev.FilePath != "/f.go" {
		t.Fatalf("FilePath = %q", ev.FilePath)
	}
}
