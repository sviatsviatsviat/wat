package cmdast

import (
	"testing"
)

func TestDetectDialect_sniffPowerShell(t *testing.T) {
	if g := DetectDialect(`Get-ChildItem D:\Sandbox`, ""); g != DialectPowerShell {
		t.Fatalf("got %q", g)
	}
}

func TestDetectDialect_sniffBash(t *testing.T) {
	if g := DetectDialect(`cd /tmp/example && git status`, ""); g != DialectBash {
		t.Fatalf("got %q", g)
	}
}

func TestDetectDialect_hostHintPowerShell(t *testing.T) {
	if g := DetectDialect(`echo hello`, HostHintPowerShell); g != DialectPowerShell {
		t.Fatalf("got %q", g)
	}
}

func TestDetectDialect_fallbackBash(t *testing.T) {
	if g := DetectDialect(`echo hello`, ""); g != DialectBash {
		t.Fatalf("got %q", g)
	}
}
