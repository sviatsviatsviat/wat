package psconv

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cmdast"
)

// skipIfNoPowerShell skips the test when neither pwsh nor powershell is on PATH
// (same resolution order as [Parse]).
func skipIfNoPowerShell(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("pwsh"); err == nil {
		return
	}
	if _, err := exec.LookPath("powershell"); err == nil {
		return
	}
	t.Skip("neither pwsh nor powershell on PATH")
}

func TestParse_chainAndPipeline(t *testing.T) {
	skipIfNoPowerShell(t)
	raw := `cd ../other-repo && git push --force origin main 2>&1 | tee log.txt`
	cl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	// Windows PowerShell 5.1 does not support && pipeline chains; PS 7+ does.
	if cl.ParsePartial && strings.Contains(cl.ParseError, "not a valid statement separator") {
		t.Skip("pipeline chain with && requires PowerShell 7+")
	}
	if cl.ParsePartial {
		t.Fatalf("unexpected partial parse: %q", cl.ParseError)
	}
	if len(cl.Statements) != 1 {
		t.Fatalf("statements: %d", len(cl.Statements))
	}
	st := cl.Statements[0]
	if st.Kind != cmdast.StmtChain || st.Chain == nil || st.Chain.Operator != "&&" {
		t.Fatalf("root: %+v", st)
	}
	left := st.Chain.Left
	if left.Kind != cmdast.StmtCommand || left.Command == nil || left.Command.Name != "cd" {
		t.Fatalf("left cmd: %+v", left)
	}
	if len(left.Command.Args) < 1 || !strings.Contains(left.Command.Args[0].Literal, "other-repo") {
		t.Fatalf("cd args: %+v", left.Command.Args)
	}
	right := st.Chain.Right
	if right.Kind != cmdast.StmtPipeline || right.Pipeline == nil || len(right.Pipeline.Stages) != 2 {
		t.Fatalf("right pipeline: %+v", right)
	}
	g := right.Pipeline.Stages[0].Command
	if g == nil || g.Name != "git" {
		t.Fatalf("git cmd: %+v", g)
	}
	if len(g.Redirects) < 1 {
		t.Fatalf("expected redirect 2>&1 on git, got %+v", g.Redirects)
	}
	if !strings.Contains(g.Redirects[0].Operator, ">") {
		t.Fatalf("redirect op: %+v", g.Redirects[0])
	}
	if !psArgLiteralsContain(g.Args, "origin") || !psArgLiteralsContain(g.Args, "main") {
		t.Fatalf("git args missing origin/main: %+v", g.Args)
	}
	tee := right.Pipeline.Stages[1].Command
	if tee == nil || tee.Name != "tee" {
		t.Fatalf("tee: %+v", tee)
	}
	if len(tee.Args) < 1 || !strings.Contains(tee.Args[0].Literal, "log.txt") {
		t.Fatalf("tee args: %+v", tee.Args)
	}
}

func TestParse_simpleCommand(t *testing.T) {
	skipIfNoPowerShell(t)
	raw := `Write-Output hello`
	cl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if cl.Lang != cmdast.LangPowerShell {
		t.Fatalf("lang %q", cl.Lang)
	}
	if len(cl.Statements) != 1 || cl.Statements[0].Kind != cmdast.StmtCommand {
		t.Fatalf("statements: %+v", cl.Statements)
	}
	cmd := cl.Statements[0].Command
	if cmd == nil || cmd.Name != "Write-Output" {
		t.Fatalf("command: %+v", cmd)
	}
	found := false
	for _, a := range cmd.Args {
		if strings.TrimSpace(a.Literal) == "hello" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("args: %+v", cmd.Args)
	}
}

func TestParse_pipelineTwoStages(t *testing.T) {
	skipIfNoPowerShell(t)
	raw := `Write-Output 1 | ForEach-Object { $_ }`
	cl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Statements) < 1 {
		t.Fatalf("no statements: %+v", cl)
	}
	st := cl.Statements[0]
	if st.Kind != cmdast.StmtPipeline {
		t.Fatalf("want pipeline, got %+v", st)
	}
	if len(st.Pipeline.Stages) != 2 {
		t.Fatalf("stages: %d", len(st.Pipeline.Stages))
	}
	a := st.Pipeline.Stages[0].Command
	b := st.Pipeline.Stages[1].Command
	if a == nil || a.Name != "Write-Output" {
		t.Fatalf("stage0: %+v", a)
	}
	if b == nil || b.Name != "ForEach-Object" {
		t.Fatalf("stage1: %+v", b)
	}
}

func TestParse_twoStatementsSemicolon(t *testing.T) {
	skipIfNoPowerShell(t)
	raw := `Write-Output a; Write-Output b`
	cl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Statements) != 2 {
		t.Fatalf("want 2 statements, got %d: %+v", len(cl.Statements), cl.Statements)
	}
	if cl.Statements[0].Command == nil || cl.Statements[0].Command.Name != "Write-Output" {
		t.Fatalf("stmt0: %+v", cl.Statements[0])
	}
	if cl.Statements[1].Command == nil || cl.Statements[1].Command.Name != "Write-Output" {
		t.Fatalf("stmt1: %+v", cl.Statements[1])
	}
}

func TestParse_parseErrorPartialAST(t *testing.T) {
	skipIfNoPowerShell(t)
	raw := `(`
	cl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if !cl.ParsePartial {
		t.Fatal("expected parse_partial")
	}
	if cl.ParseError == "" {
		t.Fatal("expected parse_error")
	}
	if len(cl.Statements) < 1 || cl.Statements[0].Kind != cmdast.StmtCompound {
		t.Fatalf("expected compound statement, got %+v", cl.Statements)
	}
}

func TestParse_variableInArg(t *testing.T) {
	skipIfNoPowerShell(t)
	raw := `Write-Output $PWD`
	cl, err := Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	args := cl.Statements[0].Command.Args
	if len(args) != 1 {
		t.Fatalf("args: %+v", args)
	}
	if len(args[0].Vars) < 1 || args[0].Vars[0].Name != "PWD" {
		t.Fatalf("vars: %+v", args[0])
	}
	if !args[0].Expanded {
		t.Fatal("expected expanded for variable arg")
	}
}

func TestResolvePowerShell_errorsWhenMissing(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	_, err := resolvePowerShell()
	if err == nil {
		t.Fatal("expected error when PATH has no pwsh or powershell")
	}
}

// psArgLiteralsContain reports whether s appears in argument literals or flag-bound values (PowerShell).
func psArgLiteralsContain(args []cmdast.Arg, s string) bool {
	for _, a := range args {
		if strings.Contains(a.Literal, s) {
			return true
		}
		if a.Flag != nil && a.Flag.Value != nil && strings.Contains(a.Flag.Value.Literal, s) {
			return true
		}
	}
	return false
}
