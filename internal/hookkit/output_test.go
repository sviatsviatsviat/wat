package hookkit

import "testing"

func TestValidateEncodePair(t *testing.T) {
	t.Parallel()
	out := struct{}{}
	if err := ValidateEncodePair("claude", "preToolUse", out, []string{"preToolUse"}, nil); err != nil {
		t.Fatal(err)
	}
	if err := ValidateEncodePair("claude", "stop", out, []string{"preToolUse"}, nil); err == nil {
		t.Fatal("expected incompatible event error")
	}
	canonicalize := func(name string) (string, bool) {
		if name == "PreToolUse" {
			return "PreToolUse", true
		}
		return name, false
	}
	if err := ValidateEncodePair("copilot", "PreToolUse", out, []string{"PreToolUse"}, canonicalize); err != nil {
		t.Fatal(err)
	}
}

func TestNormalizeOutput(t *testing.T) {
	t.Parallel()
	type out struct {
		X int
	}
	v := out{X: 1}
	if got := NormalizeOutput(&v); got.(out).X != 1 {
		t.Fatal("NormalizeOutput should dereference pointer")
	}
	if NormalizeOutput(nil) != nil {
		t.Fatal("NormalizeOutput(nil) should be nil")
	}
}

func TestIsZeroOutput(t *testing.T) {
	t.Parallel()
	type zeroable struct{}
	if !IsZeroOutput(zeroable{}) {
		t.Fatal("IsZeroOutput should treat empty struct as zero")
	}
	if !IsZeroOutput(nil) {
		t.Fatal("IsZeroOutput(nil) should be true")
	}
}
