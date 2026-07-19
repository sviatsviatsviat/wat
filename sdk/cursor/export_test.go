package cursor

// Test-only Results accessors for external tests (cursor_test).
// Production code must use On*-injected Results; these are not part of the public API surface in non-test builds.

func PermissionResultsForTest() PermissionResults { return permissionResults{} }

func PostToolResultsForTest() PostToolResults { return postToolResults{} }

func StopResultsForTest() StopResults { return stopResults{} }

func SessionStartResultsForTest() SessionStartResults { return sessionStartResults{} }

func BeforeSubmitPromptResultsForTest() BeforeSubmitPromptResults {
	return beforeSubmitPromptResults{}
}

func PreCompactResultsForTest() PreCompactResults { return preCompactResults{} }

// Decode is a test-only alias for codec.Decode (available to cursor_test).
func Decode(raw []byte) (Event, error) {
	ev, err := codec.Decode(raw)
	if err != nil {
		return nil, err
	}
	if ev == nil {
		return nil, nil
	}
	return ev.(Event), nil
}
