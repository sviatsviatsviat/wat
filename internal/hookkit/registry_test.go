package hookkit

import "testing"

func TestDecoderRegistry(t *testing.T) {
	t.Parallel()
	r := NewDecoderRegistry()
	if _, ok := r.Lookup("PreToolUse"); ok {
		t.Fatal("expected empty registry")
	}

	called := false
	r.Register("PreToolUse", func(raw []byte, received, canonical string) (any, error) {
		called = true
		if string(raw) != `{}` || received != "pre" || canonical != "PreToolUse" {
			t.Fatalf("unexpected args: raw=%q received=%q canonical=%q", raw, received, canonical)
		}
		return "ok", nil
	})

	fn, ok := r.Lookup("PreToolUse")
	if !ok {
		t.Fatal("missing decoder")
	}
	got, err := fn([]byte(`{}`), "pre", "PreToolUse")
	if err != nil {
		t.Fatal(err)
	}
	if got != "ok" || !called {
		t.Fatalf("got = %v called = %v", got, called)
	}

	other := NewDecoderRegistry()
	if _, ok := other.Lookup("PreToolUse"); ok {
		t.Fatal("registries must be isolated")
	}
}
