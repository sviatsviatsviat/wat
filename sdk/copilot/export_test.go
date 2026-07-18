package copilot

// Test-only Results accessors for external tests (copilot_test).

func PreToolResultsForTest() PreToolResults { return preToolResults{} }

func PostToolResultsForTest() PostToolResults { return postToolResults{} }

func StopResultsForTest() StopResults { return stopResults{} }

func PermissionRequestResultsForTest() PermissionRequestResults {
	return permissionRequestResults{}
}

func PostToolFailureResultsForTest() PostToolFailureResults {
	return postToolFailureResults{}
}

func SessionStartResultsForTest() SessionStartResults { return sessionStartResults{} }

func SubagentStartResultsForTest() SubagentStartResults { return subagentStartResults{} }

func NotificationResultsForTest() NotificationResults { return notificationResults{} }

// Decode is a test-only alias for decode (available to copilot_test).
func Decode(raw []byte) (Event, error) {
	ev, err := decode(raw)
	if err != nil {
		return nil, err
	}
	if ev == nil {
		return nil, nil
	}
	return ev.(Event), nil
}
