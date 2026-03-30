package shellparse

import (
	"strings"

	"github.com/sviatsviatsviat/wat/internal/core"
	"mvdan.cc/sh/v3/syntax"
)

// BashParser parses POSIX/Bash-like command strings using [syntax.Parser] from mvdan.cc/sh/v3.
//
// It walks top-level statements and compound commands that are common in one-liners:
// simple calls, && / || chains, pipelines (| and |&), subshells, and braced blocks.
// Commands inside other constructs (if, while, for, case, etc.) are not expanded—those
// nodes are skipped so the result focuses on agent-style shell invocations.
//
// Words are turned into strings with a small literal/quoting helper ([wordToString]);
// parameter expansion, arithmetic, globs, and command substitutions are not evaluated.
type BashParser struct{}

// Dialect returns [core.DialectBash].
func (p *BashParser) Dialect() string { return core.DialectBash }

// Parse builds a [core.ParseResult] with Raw set to raw. On syntax error it returns
// the partial result (dialect and raw only) and the parse error from the shell parser.
// On success, Pipeline lists [core.CommandNode] values in execution order for the
// constructs described on [BashParser].
func (p *BashParser) Parse(raw string) (core.ParseResult, error) {
	result := core.ParseResult{Dialect: core.DialectBash, Raw: raw}
	prog, err := syntax.NewParser().Parse(strings.NewReader(raw), "")
	if err != nil {
		return result, err
	}
	for _, stmt := range prog.Stmts {
		result.Pipeline = append(result.Pipeline, collectFromStmt(stmt)...)
	}
	return result, nil
}

// collectFromStmt extracts command nodes from a parsed statement, handling calls,
// &&/||, pipes, subshells, and blocks as described on [BashParser].
func collectFromStmt(stmt *syntax.Stmt) []core.CommandNode {
	if stmt == nil || stmt.Cmd == nil {
		return nil
	}
	switch cmd := stmt.Cmd.(type) {
	case *syntax.CallExpr:
		n := callExprToNode(cmd, 0, 1)
		return []core.CommandNode{n}
	case *syntax.BinaryCmd:
		switch cmd.Op {
		case syntax.AndStmt, syntax.OrStmt:
			return append(collectFromStmt(cmd.X), collectFromStmt(cmd.Y)...)
		case syntax.Pipe, syntax.PipeAll:
			parts := flattenPipe(stmt)
			plen := len(parts)
			var out []core.CommandNode
			for i, st := range parts {
				if st == nil || st.Cmd == nil {
					continue
				}
				ce, ok := st.Cmd.(*syntax.CallExpr)
				if !ok {
					continue
				}
				n := callExprToNode(ce, i, plen)
				out = append(out, n)
			}
			return out
		}
	case *syntax.Subshell:
		var out []core.CommandNode
		for _, s := range cmd.Stmts {
			out = append(out, collectFromStmt(s)...)
		}
		return out
	case *syntax.Block:
		var out []core.CommandNode
		for _, s := range cmd.Stmts {
			out = append(out, collectFromStmt(s)...)
		}
		return out
	}
	return nil
}

// flattenPipe returns the pipeline segment statements in order for a left-associative
// pipe tree (a | b | c).
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

// callExprToNode maps a simple command’s words to a [core.CommandNode]. Arguments starting with "-"
// are recorded as switches; other non-name tokens go to Args. Flags as separate name=value
// pairs are not distinguished from POSIX long options here.
func callExprToNode(call *syntax.CallExpr, pipeIndex, pipeLen int) core.CommandNode {
	cmd := core.CommandNode{
		Flags:      make(map[string]string),
		PipeIndex:  pipeIndex,
		PipeLength: pipeLen,
	}
	if len(call.Args) == 0 {
		return cmd
	}
	cmd.Name = wordToString(call.Args[0])
	for _, arg := range call.Args[1:] {
		val := wordToString(arg)
		if strings.HasPrefix(val, switchPrefix) {
			cmd.Switches = append(cmd.Switches, val)
		} else {
			cmd.Args = append(cmd.Args, val)
		}
	}
	return cmd
}

// wordToString recovers a best-effort string from a shell [syntax.Word]: pure literals,
// simple single- and double-quoted segments. It returns the empty string for words that
// require full shell expansion (parameters, command substitution, etc.).
func wordToString(w *syntax.Word) string {
	if w == nil {
		return ""
	}
	if s := w.Lit(); s != "" {
		return s
	}
	var b strings.Builder
	for _, p := range w.Parts {
		switch x := p.(type) {
		case *syntax.Lit:
			b.WriteString(x.Value)
		case *syntax.SglQuoted:
			b.WriteString(x.Value)
		case *syntax.DblQuoted:
			for _, q := range x.Parts {
				if l, ok := q.(*syntax.Lit); ok {
					b.WriteString(l.Value)
				}
			}
		}
	}
	return b.String()
}
