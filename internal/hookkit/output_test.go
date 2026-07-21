package hookkit

import "testing"

type testOutput struct {
	zero    bool
	allowed []string
}

func (o testOutput) IsZero() bool            { return o.zero }
func (o testOutput) AllowedEvents() []string { return o.allowed }
func (o testOutput) Encode(string) ([]byte, int, error) {
	return nil, 0, nil
}

func TestValidateEncodePair(t *testing.T) {
	t.Parallel()
	out := testOutput{allowed: []string{"preToolUse"}}
	if err := ValidateEncodePair("claude", "preToolUse", out, nil); err != nil {
		t.Fatal(err)
	}
	if err := ValidateEncodePair("claude", "stop", out, nil); err == nil {
		t.Fatal("expected incompatible event error")
	}
	canonicalize := func(name string) (string, bool) {
		if name == "PreToolUse" {
			return "PreToolUse", true
		}
		return name, false
	}
	out = testOutput{allowed: []string{"PreToolUse"}}
	if err := ValidateEncodePair("copilot", "PreToolUse", out, canonicalize); err != nil {
		t.Fatal(err)
	}
}
