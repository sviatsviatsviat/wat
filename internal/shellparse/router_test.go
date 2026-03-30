package shellparse

import (
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

func TestParserRouter_Parse_bashMarker(t *testing.T) {
	r := NewParserRouter()
	res, err := r.Parse(`sudo /usr/bin/git log -1`)
	if err != nil {
		t.Fatal(err)
	}
	if res.Dialect != core.DialectBash {
		t.Fatalf("got %q", res.Dialect)
	}
}

func TestParserRouter_Parse_powershellMarker(t *testing.T) {
	r := NewParserRouter()
	res, err := r.Parse(`Write-Output $env:TEMP`)
	if err != nil {
		t.Fatal(err)
	}
	if res.Dialect != core.DialectPowerShell {
		t.Fatalf("got %q", res.Dialect)
	}
}
