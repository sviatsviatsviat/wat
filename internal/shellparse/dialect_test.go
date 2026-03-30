package shellparse

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func TestDetectDialect_sniffPowerShell(t *testing.T) {
	if g := detectDialect(`Get-ChildItem D:\Sandbox`, ""); g != core.DialectPowerShell {
		t.Fatalf("got %q", g)
	}
}

func TestDetectDialect_sniffBash(t *testing.T) {
	if g := detectDialect(`cd /tmp/example && git status`, ""); g != core.DialectBash {
		t.Fatalf("got %q", g)
	}
}

func TestDetectDialect_hostHintPowerShell(t *testing.T) {
	// No markers: plain token; Windows-oriented hint selects PowerShell.
	if g := detectDialect(`echo hello`, hostHintPowerShell); g != core.DialectPowerShell {
		t.Fatalf("got %q", g)
	}
}

func TestDetectDialect_fallbackBash(t *testing.T) {
	if g := detectDialect(`echo hello`, ""); g != core.DialectBash {
		t.Fatalf("got %q", g)
	}
}
