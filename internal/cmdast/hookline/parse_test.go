package hookline

import (
	"os/exec"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cmdast"
)

func TestParseCommandLine_bashMarkers(t *testing.T) {
	raw := `true && echo ok`
	cl, dialect, err := ParseCommandLine(raw)
	if err != nil {
		t.Fatal(err)
	}
	if dialect != cmdast.DialectBash {
		t.Fatalf("dialect %q", dialect)
	}
	if len(cl.Statements) != 1 || cl.Statements[0].Kind != cmdast.StmtChain {
		t.Fatalf("root: %+v", cl.Statements)
	}
}

func TestParseCommandLine_chainShape(t *testing.T) {
	raw := `git status && echo x`
	cl, _, err := ParseCommandLine(raw)
	if err != nil {
		t.Fatal(err)
	}
	st := cl.Statements[0]
	if st.Kind != cmdast.StmtChain || st.Chain == nil || st.Chain.Operator != "&&" {
		t.Fatalf("want && chain: %+v", st)
	}
	if st.Chain.Left == nil || st.Chain.Left.Command == nil || st.Chain.Left.Command.Name != "git" {
		t.Fatalf("left: %+v", st.Chain.Left)
	}
	if st.Chain.Right == nil || st.Chain.Right.Command == nil || st.Chain.Right.Command.Name != "echo" {
		t.Fatalf("right: %+v", st.Chain.Right)
	}
}

func TestParseCommandLine_sudoBashMarker(t *testing.T) {
	cl, d, err := ParseCommandLine(`sudo /usr/bin/git log -1`)
	if err != nil {
		t.Fatal(err)
	}
	if d != cmdast.DialectBash {
		t.Fatalf("dialect %q", d)
	}
	if cl == nil || len(cl.Statements) < 1 {
		t.Fatalf("command line: %+v", cl)
	}
}

func TestParseCommandLine_powershellMarker(t *testing.T) {
	if _, err := exec.LookPath("pwsh"); err != nil {
		t.Skip("pwsh not on PATH")
	}
	cl, d, err := ParseCommandLine(`Write-Output $env:TEMP`)
	if err != nil {
		t.Fatal(err)
	}
	if d != cmdast.DialectPowerShell {
		t.Fatalf("dialect %q", d)
	}
	if cl == nil || cl.Lang != cmdast.LangPowerShell {
		t.Fatalf("got %+v", cl)
	}
}
