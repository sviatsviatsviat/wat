package shconv

import (
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cmdast"
	"mvdan.cc/sh/v3/syntax"
)

func convertStmt(stmt *syntax.Stmt, raw string) *cmdast.Statement {
	if stmt == nil || stmt.Cmd == nil {
		return nil
	}
	bg := stmt.Background
	neg := stmt.Negated
	redirs := convertRedirects(stmt.Redirs, raw)

	switch cmd := stmt.Cmd.(type) {
	case *syntax.CallExpr:
		c := convertCallExpr(cmd, raw)
		c.Background = bg
		c.Negated = neg
		mergeCmdRedirs(c, redirs)
		return &cmdast.Statement{
			Kind:    cmdast.StmtCommand,
			Span:    spanStmt(stmt, raw),
			Command: c,
		}

	case *syntax.BinaryCmd:
		switch cmd.Op {
		case syntax.Pipe, syntax.PipeAll:
			return buildPipeline(stmt, raw, bg)
		case syntax.AndStmt, syntax.OrStmt:
			left := convertStmt(cmd.X, raw)
			right := convertStmt(cmd.Y, raw)
			if left == nil || right == nil {
				return nil
			}
			return &cmdast.Statement{
				Kind: cmdast.StmtChain,
				Span: spanStmt(stmt, raw),
				Chain: &cmdast.Chain{
					Operator: cmd.Op.String(),
					Left:     left,
					Right:    right,
				},
			}
		default:
			return nil
		}

	case *syntax.IfClause:
		return compoundStmt(stmt, raw, cmdast.CompoundIf, cmd, neg, bg)
	case *syntax.WhileClause:
		kind := cmdast.CompoundWhile
		if cmd.Until {
			kind = cmdast.CompoundUntil
		}
		return compoundStmt(stmt, raw, kind, cmd, neg, bg)
	case *syntax.ForClause:
		kind := cmdast.CompoundFor
		switch cmd.Loop.(type) {
		case *syntax.WordIter:
			kind = cmdast.CompoundForEach
		case *syntax.CStyleLoop:
			kind = cmdast.CompoundFor
		}
		return compoundStmt(stmt, raw, kind, cmd, neg, bg)
	case *syntax.CaseClause:
		return compoundStmt(stmt, raw, cmdast.CompoundCase, cmd, neg, bg)
	case *syntax.FuncDecl:
		return compoundStmt(stmt, raw, cmdast.CompoundFunction, cmd, neg, bg)
	case *syntax.Subshell:
		return compoundStmt(stmt, raw, cmdast.CompoundSubshell, cmd, neg, bg)
	case *syntax.Block:
		return compoundStmt(stmt, raw, cmdast.CompoundBlock, cmd, neg, bg)
	case *syntax.TestClause:
		return compoundStmt(stmt, raw, cmdast.CompoundTest, cmd, neg, bg)
	case *syntax.ArithmCmd:
		return compoundStmt(stmt, raw, cmdast.CompoundArithmetic, cmd, neg, bg)
	case *syntax.LetClause:
		return compoundStmt(stmt, raw, cmdast.CompoundArithmetic, cmd, neg, bg)
	case *syntax.DeclClause:
		return declClauseStmt(stmt, cmd, raw, neg, bg, redirs)
	case *syntax.TimeClause:
		return compoundStmt(stmt, raw, cmdast.CompoundOther, cmd, neg, bg)
	case *syntax.CoprocClause:
		return compoundStmt(stmt, raw, cmdast.CompoundOther, cmd, neg, bg)
	default:
		return compoundStmt(stmt, raw, cmdast.CompoundOther, cmd, neg, bg)
	}
}

func compoundStmt(stmt *syntax.Stmt, raw string, kind cmdast.CompoundKind, cmd syntax.Node, _, _ bool) *cmdast.Statement {
	text := sliceRaw(raw, cmd.Pos(), cmd.End())
	return &cmdast.Statement{
		Kind: cmdast.StmtCompound,
		Span: spanStmt(stmt, raw),
		Compound: &cmdast.Compound{
			CompoundKind: kind,
			Raw:          text,
		},
	}
}

func declClauseStmt(stmt *syntax.Stmt, d *syntax.DeclClause, raw string, neg, bg bool, redirs []cmdast.Redirect) *cmdast.Statement {
	name := ""
	if d.Variant != nil {
		name = d.Variant.Value
	}
	c := &cmdast.Command{
		Name:        name,
		Background:  stmt.Background,
		Negated:     stmt.Negated || neg,
		Assignments: make([]cmdast.Assign, 0, len(d.Args)),
		Redirects:   redirs,
	}
	for _, a := range d.Args {
		c.Assignments = append(c.Assignments, convertShellAssign(a, raw))
	}
	return &cmdast.Statement{
		Kind:    cmdast.StmtCommand,
		Span:    spanStmt(stmt, raw),
		Command: c,
	}
}

func convertCallExpr(call *syntax.CallExpr, raw string) *cmdast.Command {
	c := &cmdast.Command{}
	if len(call.Args) == 0 {
		return c
	}
	c.Name = commandName(call.Args[0], raw)
	args := make([]cmdast.Arg, 0, len(call.Args)-1)
	for _, w := range call.Args[1:] {
		args = append(args, convertWord(w, raw))
	}
	args = linkFlagValues(args)
	c.Args = args
	for _, a := range call.Assigns {
		c.Assignments = append(c.Assignments, convertShellAssign(a, raw))
	}
	return c
}

func mergeCmdRedirs(c *cmdast.Command, redirs []cmdast.Redirect) {
	if len(redirs) == 0 {
		return
	}
	c.Redirects = append(c.Redirects, redirs...)
}

func convertShellAssign(a *syntax.Assign, raw string) cmdast.Assign {
	op := "="
	if a.Append {
		op = "+="
	}
	name := ""
	if a.Name != nil {
		name = a.Name.Value
	}
	val := cmdast.Arg{}
	if a.Value != nil {
		val = convertWord(a.Value, raw)
	}
	return cmdast.Assign{
		Span:     spanNode(a, raw),
		Name:     name,
		Operator: op,
		Value:    val,
	}
}

func buildPipeline(stmt *syntax.Stmt, raw string, bg bool) *cmdast.Statement {
	parts := flattenPipe(stmt)
	stages := make([]*cmdast.Statement, 0, len(parts))
	for _, p := range parts {
		if st := convertStmt(p, raw); st != nil {
			stages = append(stages, st)
		}
	}
	if len(stages) == 0 {
		return nil
	}
	return &cmdast.Statement{
		Kind: cmdast.StmtPipeline,
		Span: spanStmt(stmt, raw),
		Pipeline: &cmdast.Pipeline{
			Stages:     stages,
			Background: bg,
		},
	}
}

func flattenPipe(stmt *syntax.Stmt) []*syntax.Stmt {
	if stmt == nil || stmt.Cmd == nil {
		return nil
	}
	b, ok := stmt.Cmd.(*syntax.BinaryCmd)
	if !ok || (b.Op != syntax.Pipe && b.Op != syntax.PipeAll) {
		return []*syntax.Stmt{stmt}
	}
	return append(flattenPipe(b.X), flattenPipe(b.Y)...)
}

func convertRedirects(rs []*syntax.Redirect, raw string) []cmdast.Redirect {
	out := make([]cmdast.Redirect, 0, len(rs))
	for _, r := range rs {
		if r == nil {
			continue
		}
		op := r.Op.String()
		target := ""
		if r.Word != nil {
			target = wordRaw(r.Word, raw)
		} else if r.Hdoc != nil {
			target = wordRaw(r.Hdoc, raw)
		}
		out = append(out, cmdast.Redirect{
			Span:     spanNode(r, raw),
			Operator: op,
			Target:   target,
		})
	}
	return out
}

func convertWord(w *syntax.Word, raw string) cmdast.Arg {
	if w == nil {
		return cmdast.Arg{}
	}
	a := cmdast.Arg{
		Span:    spanNode(w, raw),
		Literal: wordRaw(w, raw),
	}
	expanded := false
	for _, p := range w.Parts {
		mergeWordPart(&a, p, raw, &expanded)
	}
	if expanded {
		a.Expanded = true
	}
	detectPOSIXFlag(&a)
	return a
}

func mergeWordPart(a *cmdast.Arg, p syntax.WordPart, raw string, expanded *bool) {
	switch x := p.(type) {
	case *syntax.Lit:
		// literal only
	case *syntax.SglQuoted:
		// literal only
	case *syntax.DblQuoted:
		*expanded = true
		for _, q := range x.Parts {
			mergeWordPart(a, q, raw, expanded)
		}
	case *syntax.ParamExp:
		a.Vars = append(a.Vars, varFromParam(x, raw))
	case *syntax.CmdSubst:
		a.CmdSubs = append(a.CmdSubs, convertCmdSubst(x, raw))
		*expanded = true
	case *syntax.ProcSubst:
		a.CmdSubs = append(a.CmdSubs, convertProcSubst(x, raw))
		*expanded = true
	case *syntax.ArithmExp:
		*expanded = true
	case *syntax.ExtGlob, *syntax.BraceExp:
		*expanded = true
	default:
		*expanded = true
	}
}

func varFromParam(p *syntax.ParamExp, raw string) cmdast.VarRef {
	name := ""
	if p.Param != nil {
		name = p.Param.Value
	}
	return cmdast.VarRef{
		Span: spanNode(p, raw),
		Name: name,
	}
}

func convertCmdSubst(s *syntax.CmdSubst, raw string) cmdast.CmdSub {
	sts := make([]*cmdast.Statement, 0, len(s.Stmts))
	for _, st := range s.Stmts {
		if c := convertStmt(st, raw); c != nil {
			sts = append(sts, c)
		}
	}
	return cmdast.CmdSub{
		Span:     spanNode(s, raw),
		Commands: sts,
	}
}

func convertProcSubst(p *syntax.ProcSubst, raw string) cmdast.CmdSub {
	sts := make([]*cmdast.Statement, 0, len(p.Stmts))
	for _, st := range p.Stmts {
		if c := convertStmt(st, raw); c != nil {
			sts = append(sts, c)
		}
	}
	return cmdast.CmdSub{
		Span:     spanNode(p, raw),
		Commands: sts,
	}
}

func commandName(w *syntax.Word, raw string) string {
	s := strings.TrimSpace(wordRaw(w, raw))
	if len(s) >= 2 && ((s[0] == '\'' && s[len(s)-1] == '\'') || (s[0] == '"' && s[len(s)-1] == '"')) {
		return s[1 : len(s)-1]
	}
	return s
}

func wordRaw(w *syntax.Word, raw string) string {
	if w == nil {
		return ""
	}
	return sliceRaw(raw, w.Pos(), w.End())
}

func sliceRaw(raw string, pos, end syntax.Pos) string {
	a, b := int(pos.Offset()), int(end.Offset())
	if a < 0 || b > len(raw) || a > b {
		return ""
	}
	return raw[a:b]
}

func spanStmt(stmt *syntax.Stmt, raw string) cmdast.Span {
	return spanNode(stmt, raw)
}

func spanNode(n syntax.Node, raw string) cmdast.Span {
	if n == nil {
		return cmdast.Span{}
	}
	p := n.Pos()
	e := n.End()
	return cmdast.Span{
		Start: cmdast.Position{Offset: uint(p.Offset()), Line: uint(p.Line()), Col: uint(p.Col())},
		End:   cmdast.Position{Offset: uint(e.Offset()), Line: uint(e.Line()), Col: uint(e.Col())},
	}
}

func detectPOSIXFlag(a *cmdast.Arg) {
	lit := strings.TrimSpace(a.Literal)
	if a.Flag != nil {
		return
	}
	if strings.HasPrefix(lit, "--") && len(lit) > 2 {
		body := lit[2:]
		if idx := strings.IndexByte(body, '='); idx >= 0 {
			body = body[:idx]
		}
		a.Flag = &cmdast.Flag{Name: body, Dashes: 2}
		return
	}
	if strings.HasPrefix(lit, "-") && len(lit) > 1 && lit[1] != '-' {
		name := lit[1:]
		if idx := strings.IndexByte(name, '='); idx >= 0 {
			name = name[:idx]
		}
		a.Flag = &cmdast.Flag{Name: name, Dashes: 1}
	}
}

func linkFlagValues(args []cmdast.Arg) []cmdast.Arg {
	if len(args) == 0 {
		return args
	}
	out := make([]cmdast.Arg, 0, len(args))
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a.Flag != nil && a.Flag.Value == nil && i+1 < len(args) {
			next := args[i+1]
			if a.Flag.Name != "" && !looksLikeFlag(next.Literal) {
				v := next
				a.Flag.Value = &v
				i++
				out = append(out, a)
				continue
			}
		}
		out = append(out, a)
	}
	return out
}

func looksLikeFlag(s string) bool {
	t := strings.TrimSpace(s)
	return strings.HasPrefix(t, "-") && len(t) > 1
}
