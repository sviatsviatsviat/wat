package shellparse

import (
	"os/exec"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// Sanitized PowerShell-style invocations (paths and identities replaced).
var psParseFixtures = []struct {
	name string
	raw  string
}{
	{
		name: "simple_remove",
		raw:  `Remove-Item -Recurse -Force 'D:\Sandbox\project'`,
	},
	{
		name: "pipeline_select",
		raw:  `Get-ChildItem 'C:\Contoso\Repo' | Select-Object Name`,
	},
	{
		name: "invoke_expression_obfuscated",
		raw:  `powershell -NoProfile -File 'D:\Contoso\Repo\scripts\continual-learning\delta.ps1' -TranscriptsRoot 'C:\Users\Public\agent-transcripts'`,
	},
}

func TestPowerShellParser_parseFast_useCases(t *testing.T) {
	p := &PowerShellParser{}
	for _, tc := range psParseFixtures {
		t.Run(tc.name, func(t *testing.T) {
			res, err := p.Parse(tc.raw)
			if err != nil {
				t.Fatal(err)
			}
			if res.Dialect != core.DialectPowerShell {
				t.Fatalf("dialect %q", res.Dialect)
			}
			if len(res.Pipeline) < 1 {
				t.Fatalf("empty pipeline")
			}
			if res.Pipeline[0].Name == "" {
				t.Fatal("empty command name")
			}
		})
	}
}

func TestPowerShellParser_complex_pwshIfAvailable(t *testing.T) {
	if _, err := exec.LookPath("pwsh"); err != nil {
		t.Skip("pwsh not on PATH")
	}
	// Complexity trigger: subexpression (sanitized path only).
	raw := `$(Get-Location).Path + '\example'`
	p := &PowerShellParser{}
	res, err := p.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if res.Dialect != core.DialectPowerShell {
		t.Fatalf("dialect %q", res.Dialect)
	}
	if len(res.Pipeline) < 1 {
		t.Fatalf("expected at least one node from pwsh or fast fallback: %+v", res.Pipeline)
	}
}
