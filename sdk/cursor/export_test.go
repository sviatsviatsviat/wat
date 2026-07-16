package cursor

// Test-only Results accessors for external tests (cursor_test).
// Production code must use Chain-injected Results; these are not part of the public API surface in non-test builds.

func PermissionResultsForTest() PermissionResults { return permissionResults{} }

func PostToolResultsForTest() PostToolResults { return postToolResults{} }

func StopResultsForTest() StopResults { return stopResults{} }

func SessionStartResultsForTest() SessionStartResults { return sessionStartResults{} }

func BeforeSubmitPromptResultsForTest() BeforeSubmitPromptResults {
	return beforeSubmitPromptResults{}
}

func PreCompactResultsForTest() PreCompactResults { return preCompactResults{} }
