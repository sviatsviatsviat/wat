package claude

import (
	"testing"
)

func TestDecode_Setup(t *testing.T) {
	mustDecode[Setup](t, `{"session_id":"s","hook_event_name":"Setup","trigger":"init"}`, EventSetup)
}
