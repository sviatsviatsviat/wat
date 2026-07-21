package hookkit

import (
	"maps"
	"testing"
)

func TestTakeLastPtr(t *testing.T) {
	a, b := 1, 2
	got, w := TakeLastPtr("x", &a, &b)
	if got == nil || *got != 2 || w == "" {
		t.Fatalf("got=%v warn=%q", got, w)
	}
	got, w = TakeLastPtr("x", &a, (*int)(nil))
	if got == nil || *got != 1 || w != "" {
		t.Fatalf("dst only: got=%v warn=%q", got, w)
	}
}

func TestTakeLastString(t *testing.T) {
	got, w := TakeLastString("field", "a", "b")
	if got != "b" || w != OverwriteWarning("field") {
		t.Fatalf("got=%q warn=%q", got, w)
	}
	got, w = TakeLastString("field", "a", "")
	if got != "a" || w != "" {
		t.Fatalf("empty src: got=%q warn=%q", got, w)
	}
}

func TestTakeLastMap(t *testing.T) {
	dst := map[string]int{"a": 1}
	src := map[string]int{"b": 2}
	orig := maps.Clone(dst)
	got, w := TakeLastMap("env", dst, src)
	if w == "" || got["b"] != 2 {
		t.Fatalf("got=%v warn=%q", got, w)
	}
	if !maps.Equal(dst, orig) {
		t.Fatalf("caller map mutated: %v", dst)
	}
	got["b"] = 99
	if src["b"] != 2 {
		t.Fatalf("src aliased")
	}
}

func TestMergeRankedString(t *testing.T) {
	d, r := MergeRankedString("allow", "ok", "deny", "blocked", PermissionRankString)
	if d != "deny" || r != "blocked" {
		t.Fatalf("got %q %q", d, r)
	}
	d, r = MergeRankedString("deny", "blocked", "allow", "ok", PermissionRankString)
	if d != "deny" || r != "blocked" {
		t.Fatalf("weaker ignored: %q %q", d, r)
	}
	d, r = MergeRankedString("allow", "ok", "deny", "", PermissionRankString)
	if d != "deny" || r != "" {
		t.Fatalf("clear reason: %q %q", d, r)
	}
}

func TestJoinContextStrings(t *testing.T) {
	if got := JoinContextStrings("a", "b"); got != "a\n\nb" {
		t.Fatalf("got %q", got)
	}
}
