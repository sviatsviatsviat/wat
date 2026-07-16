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
