package tools

import (
	"encoding/json"
	"testing"
)

func TestInputAsBash(t *testing.T) {
	in := NewInput(ToolBash, "Bash", json.RawMessage(`{"command":"go test ./..."}`))
	got, ok := in.AsBash()
	if !ok || got.Command != "go test ./..." {
		t.Fatalf("AsBash = %+v, %v", got, ok)
	}
	if _, ok := in.AsWrite(); ok {
		t.Fatal("AsWrite should be false for bash input")
	}
}

func TestInputAsWritePathAlias(t *testing.T) {
	in := NewInput(ToolWrite, "Write", json.RawMessage(`{"file_path":"/a","content":"x"}`))
	got, ok := in.AsWrite()
	if !ok || got.Path != "/a" || got.Content != "x" {
		t.Fatalf("AsWrite = %+v, %v", got, ok)
	}
}
