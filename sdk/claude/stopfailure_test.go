package claude

import (
	"testing"
)

func TestDecode_StopFailure(t *testing.T) {
	mustDecode[StopFailure](t, `{"session_id":"s","hook_event_name":"StopFailure","error_type":"rate_limit","message":"slow down"}`, EventStopFailure)
}
