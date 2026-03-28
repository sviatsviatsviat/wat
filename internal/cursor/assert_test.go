package cursor

import "testing"

// cursorHookStdoutSuccessLine is the hook stdout line successful Cursor hooks must write.
// Tests use this value instead of the production default response line constant so
// assertions stay explicit if that symbol is renamed or the contract is revisited.
const cursorHookStdoutSuccessLine = "{}\n"

func assertEqual(t *testing.T, want, got string) {
	t.Helper()
	if want != got {
		t.Fatalf("want %q, got %q", want, got)
	}
}

func assertIntEqual(t *testing.T, want, got int) {
	t.Helper()
	if want != got {
		t.Fatalf("want %d, got %d", want, got)
	}
}
