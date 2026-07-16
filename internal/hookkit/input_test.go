package hookkit

import (
	"encoding/json"
	"testing"
)

func TestInput_As(t *testing.T) {
	t.Parallel()
	type bash struct {
		Command string `json:"command"`
	}
	in := NewInput("Bash", json.RawMessage(`{"command":"ls"}`))
	got, ok := As[bash](in, "Bash")
	if !ok || got.Command != "ls" {
		t.Fatalf("As = %+v, %v", got, ok)
	}
	if _, ok := As[bash](in, "Write"); ok {
		t.Fatal("As should fail for wrong tool name")
	}
	if _, ok := As[bash](in, "bash"); ok {
		t.Fatal("As should be case-sensitive")
	}
}

func TestInput_AsFold(t *testing.T) {
	t.Parallel()
	type bash struct {
		Command string `json:"command"`
	}
	in := NewInput("PowerShell", json.RawMessage(`{"command":"Get-ChildItem"}`))
	got, ok := AsFold[bash](in, "bash", "powershell", "shell")
	if !ok || got.Command != "Get-ChildItem" {
		t.Fatalf("AsFold = %+v, %v", got, ok)
	}
	if _, ok := AsFold[bash](in, "write"); ok {
		t.Fatal("AsFold should fail for unmatched names")
	}
	empty := NewInput("bash", nil)
	if got, ok := AsFold[bash](empty, "Bash"); !ok || got.Command != "" {
		t.Fatalf("AsFold empty raw = %+v, %v", got, ok)
	}
}

func TestNewInput_panicsWhenEmpty(t *testing.T) {
	t.Parallel()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()
	_ = NewInput("", nil)
}
