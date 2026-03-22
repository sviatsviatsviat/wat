package helpers

import "testing"

func TestJoinSemicolonSeparatedStrings(t *testing.T) {
	t.Parallel()
	if joined := JoinSemicolonSeparatedStrings(nil); joined != "" {
		t.Fatalf("nil: want empty, got %q", joined)
	}
	if joined := JoinSemicolonSeparatedStrings([]string{}); joined != "" {
		t.Fatalf("empty: want empty, got %q", joined)
	}
	if joined := JoinSemicolonSeparatedStrings([]string{"a"}); joined != "a" {
		t.Fatalf("single: want a, got %q", joined)
	}
	if joined := JoinSemicolonSeparatedStrings([]string{"a", "b"}); joined != "a;b" {
		t.Fatalf("two: want a;b, got %q", joined)
	}
	// JoinSemicolonSeparatedStrings does not escape ";"; delimiter and in-element semicolons both appear as ";".
	if joined := JoinSemicolonSeparatedStrings([]string{"a;b", "c"}); joined != "a;b;c" {
		t.Fatalf("semicolon in element: want a;b;c, got %q", joined)
	}
	if joined := JoinSemicolonSeparatedStrings([]string{"  x  ", "y"}); joined != "  x  ;y" {
		t.Fatalf("whitespace preserved: want leading/trailing spaces kept, got %q", joined)
	}
}

func TestStringFromPtr(t *testing.T) {
	t.Parallel()
	if dereferenced := StringFromPtr(nil); dereferenced != "" {
		t.Fatalf("nil: want empty, got %q", dereferenced)
	}
	sample := "x"
	if dereferenced := StringFromPtr(&sample); dereferenced != "x" {
		t.Fatalf("non-nil: want x, got %q", dereferenced)
	}
}
