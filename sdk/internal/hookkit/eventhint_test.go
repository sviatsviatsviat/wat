package hookkit

import "testing"

func TestEventHint(t *testing.T) {
	t.Parallel()
	cfg := DefaultEventHint()
	ApplyHintOptions(&cfg, WithEventHint("sessionStart"))
	if cfg.Hint != "sessionStart" {
		t.Fatalf("Hint = %q", cfg.Hint)
	}
}
