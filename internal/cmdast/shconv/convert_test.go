package shconv

import (
	"strings"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/cmdast"
)

func TestParse_chainAndPipeline(t *testing.T) {
	raw := `cd ../other-repo && git push --force origin main 2>&1 | tee log.txt`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
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
	if !argLiteralsContain(g.Args, "origin") || !argLiteralsContain(g.Args, "main") {
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

func TestParse_nestedChain(t *testing.T) {
	raw := `cd /tmp && rm -rf * || echo failed`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	st := cl.Statements[0]
	if st.Kind != cmdast.StmtChain || st.Chain.Operator != "||" {
		t.Fatalf("root op: %+v", st)
	}
	if st.Chain.Left.Kind != cmdast.StmtChain || st.Chain.Left.Chain.Operator != "&&" {
		t.Fatalf("left: %+v", st.Chain.Left)
	}
	and := st.Chain.Left
	if and.Chain.Left.Command == nil || and.Chain.Left.Command.Name != "cd" {
		t.Fatalf("cd: %+v", and.Chain.Left)
	}
	rm := and.Chain.Right
	if rm.Kind != cmdast.StmtCommand || rm.Command == nil || rm.Command.Name != "rm" {
		t.Fatalf("rm cmd: %+v", rm)
	}
	orRight := st.Chain.Right
	if orRight.Kind != cmdast.StmtCommand || orRight.Command == nil || orRight.Command.Name != "echo" {
		t.Fatalf("echo cmd: %+v", orRight)
	}
	if joinArgLiterals(orRight.Command.Args) == "" || !strings.Contains(joinArgLiterals(orRight.Command.Args), "failed") {
		t.Fatalf("echo args: %+v", orRight.Command.Args)
	}
}

func TestParse_compoundInPipeline(t *testing.T) {
	raw := `if grep -q error log.txt; then rm log.txt; fi | cat`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	if cl.Statements[0].Kind != cmdast.StmtPipeline {
		t.Fatalf("want pipeline root, got %+v", cl.Statements[0])
	}
	p := cl.Statements[0].Pipeline
	if len(p.Stages) != 2 {
		t.Fatalf("stages %d", len(p.Stages))
	}
	if p.Stages[0].Kind != cmdast.StmtCompound || p.Stages[0].Compound.CompoundKind != cmdast.CompoundIf {
		t.Fatalf("stage0: %+v", p.Stages[0])
	}
	if p.Stages[1].Kind != cmdast.StmtCommand || p.Stages[1].Command.Name != "cat" {
		t.Fatalf("stage1: %+v", p.Stages[1])
	}
}

func TestParse_semicolonFold(t *testing.T) {
	raw := `echo a; echo b`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Statements) != 1 || cl.Statements[0].Kind != cmdast.StmtChain {
		t.Fatalf("got %+v", cl.Statements)
	}
	if cl.Statements[0].Chain.Operator != ";" {
		t.Fatalf("op %q", cl.Statements[0].Chain.Operator)
	}
}

func TestParse_cmdSubst(t *testing.T) {
	raw := `echo $(git rev-parse HEAD)`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	arg := cl.Statements[0].Command.Args[0]
	if len(arg.CmdSubs) != 1 {
		t.Fatalf("cmd subs: %+v", arg)
	}
	sub := arg.CmdSubs[0].Commands[0].Command
	if sub == nil || sub.Name != "git" {
		t.Fatalf("inner: %+v", sub)
	}
}

func TestParse_negation(t *testing.T) {
	raw := `! true`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	if !cl.Statements[0].Command.Negated {
		t.Fatal("expected negated")
	}
}

func TestParse_flags(t *testing.T) {
	raw := `git push --force origin`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	args := cl.Statements[0].Command.Args
	found := false
	for _, a := range args {
		if a.Flag != nil && a.Flag.Name == "force" && a.Flag.Dashes == 2 {
			found = true
		}
	}
	if !found {
		var b strings.Builder
		for _, a := range args {
			b.WriteString(a.Literal)
			b.WriteByte(';')
		}
		t.Fatalf("args: %s", b.String())
	}
}

func TestParse_simplePipeline(t *testing.T) {
	raw := `echo a | cat`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Statements) != 1 || cl.Statements[0].Kind != cmdast.StmtPipeline {
		t.Fatalf("want one pipeline: %+v", cl.Statements)
	}
	st := cl.Statements[0].Pipeline.Stages
	if len(st) != 2 || st[0].Command.Name != "echo" || st[1].Command.Name != "cat" {
		t.Fatalf("stages: %+v", st)
	}
}

func TestParse_tripleSemicolonFold(t *testing.T) {
	raw := `echo a; echo b; echo c`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	if len(cl.Statements) != 1 || cl.Statements[0].Kind != cmdast.StmtChain {
		t.Fatalf("got %+v", cl.Statements)
	}
	// foldSemicolonChain is left-folded: ((echo a ; echo b) ; echo c)
	root := cl.Statements[0].Chain
	if root.Operator != ";" || root.Right.Kind != cmdast.StmtCommand {
		t.Fatalf("root: %+v", root)
	}
	if root.Right.Command.Name != "echo" || !strings.Contains(joinArgLiterals(root.Right.Command.Args), "c") {
		t.Fatalf("right echo c: %+v", root.Right)
	}
	inner := root.Left
	if inner.Kind != cmdast.StmtChain || inner.Chain.Operator != ";" {
		t.Fatalf("inner: %+v", inner)
	}
}

func TestParse_redirectStdout(t *testing.T) {
	raw := `echo hi >/dev/null`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	cmd := cl.Statements[0].Command
	if len(cmd.Redirects) < 1 {
		t.Fatalf("redirects: %+v", cmd)
	}
	if cmd.Redirects[0].Target == "" {
		t.Fatalf("redirect target: %+v", cmd.Redirects[0])
	}
}

func TestParse_backgroundSimple(t *testing.T) {
	raw := `true &`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	if !cl.Statements[0].Command.Background {
		t.Fatal("expected background")
	}
}

func TestParse_posixVariant(t *testing.T) {
	raw := `echo x`
	cl, err := Parse(raw, cmdast.LangPosix)
	if err != nil {
		t.Fatal(err)
	}
	if cl.Lang != cmdast.LangPosix {
		t.Fatalf("lang %q", cl.Lang)
	}
	if cl.Statements[0].Command.Name != "echo" {
		t.Fatalf("command: %+v", cl.Statements[0])
	}
}

func TestParse_invalidSyntaxReturnsErrorAndPartial(t *testing.T) {
	raw := `echo $(`
	cl, err := Parse(raw, cmdast.LangBash)
	if err == nil {
		t.Fatal("expected parse error")
	}
	if cl == nil || !cl.ParsePartial || cl.ParseError == "" {
		t.Fatalf("expected partial + parse_error, got cl=%v partial=%v err=%q", cl != nil, cl != nil && cl.ParsePartial, cl.ParseError)
	}
}

func TestParse_pipeAllPipeline(t *testing.T) {
	raw := `cmd1 |& cmd2`
	cl, err := Parse(raw, cmdast.LangBash)
	if err != nil {
		t.Fatal(err)
	}
	if cl.Statements[0].Kind != cmdast.StmtPipeline {
		t.Fatalf("want pipeline, got %+v", cl.Statements[0])
	}
}

func joinArgLiterals(args []cmdast.Arg) string {
	var b strings.Builder
	for _, a := range args {
		b.WriteString(a.Literal)
	}
	return b.String()
}

// argLiteralsContain reports whether s appears in argument literals or flag-bound values.
func argLiteralsContain(args []cmdast.Arg, s string) bool {
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
