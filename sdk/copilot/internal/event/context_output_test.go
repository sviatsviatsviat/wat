package event

import "testing"

func TestContextResult_EncodeAndMerge(t *testing.T) {
	a := ContextResult("one")
	b := ContextResult("two")
	merged, warnings, err := a.Merge(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %v", warnings)
	}
	out, exit, err := merged.Encode()
	if err != nil {
		t.Fatal(err)
	}
	if exit != 0 {
		t.Fatalf("exit = %d", exit)
	}
	want := `{"additional_context":"one\n\ntwo"}`
	if string(out) != want {
		t.Fatalf("Encode() = %s, want %s", out, want)
	}
	if !ContextResult("").IsZero() {
		t.Fatal("empty ContextResult should be zero")
	}
}
