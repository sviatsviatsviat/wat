package claude

// Test-only Results accessors for external tests (claude_test).
// Production code must use On*-injected Results; these are not part of the public API surface in non-test builds.

func PreToolUseResultsForTest() PreToolUseResults { return preToolUseResults{} }
func PermissionRequestResultsForTest() PermissionRequestResults {
	return permissionRequestResults{}
}
func PostToolUseResultsForTest() PostToolUseResults { return postToolUseResults{} }
func PostToolUseFailureResultsForTest() PostToolUseFailureResults {
	return postToolUseFailureResults{}
}
func UserPromptSubmitResultsForTest() UserPromptSubmitResults { return userPromptSubmitResults{} }
func StopResultsForTest() StopResults                         { return stopResults{} }
func SessionStartResultsForTest() SessionStartResults         { return sessionStartResults{} }
func PermissionDeniedResultsForTest() PermissionDeniedResults {
	return permissionDeniedResults{}
}
func MessageDisplayResultsForTest() MessageDisplayResults { return messageDisplayResults{} }
func ElicitationResultsForTest() ElicitationResults       { return elicitationResults{} }
func WorktreeCreateResultsForTest() WorktreeCreateResults { return worktreeCreateResults{} }
func NotificationResultsForTest() NotificationResults     { return notificationResults{} }
func PreCompactResultsForTest() PreCompactResults         { return preCompactResults{} }
func SubagentStartResultsForTest() SubagentStartResults   { return subagentStartResults{} }
func TaskCreatedResultsForTest() TaskCreatedResults       { return taskCreatedResults{} }
func TaskCompletedResultsForTest() TaskCompletedResults   { return taskCompletedResults{} }
func UserPromptExpansionResultsForTest() UserPromptExpansionResults {
	return userPromptExpansionResults{}
}

// Decode is a test-only alias for decode (available to claude_test).
func Decode(raw []byte) (Event, error) { return decode(raw) }
