package shconv

import (
	"io"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cmdast"
	"mvdan.cc/sh/v3/syntax"
)

// Parse parses shell source with mvdan.cc/sh and converts it to a [cmdast.CommandLine].
func Parse(raw string, lang cmdast.SourceLang) (*cmdast.CommandLine, error) {
	var variant syntax.LangVariant
	switch lang {
	case cmdast.LangPosix:
		variant = syntax.LangPOSIX
	default:
		variant = syntax.LangBash
	}
	parser := syntax.NewParser(syntax.Variant(variant))
	f, err := parser.Parse(strings.NewReader(raw), "")
	cl := Convert(f, raw, lang)
	if err != nil {
		if cl != nil {
			cl.ParseError = err.Error()
			cl.ParsePartial = true
		}
		return cl, err
	}
	return cl, nil
}

// ParseReader is like [Parse] but reads from r; filename is the parser name (e.g. "").
func ParseReader(r io.Reader, filename, raw string, lang cmdast.SourceLang) (*cmdast.CommandLine, error) {
	var variant syntax.LangVariant
	switch lang {
	case cmdast.LangPosix:
		variant = syntax.LangPOSIX
	default:
		variant = syntax.LangBash
	}
	parser := syntax.NewParser(syntax.Variant(variant))
	f, err := parser.Parse(r, filename)
	cl := Convert(f, raw, lang)
	if err != nil {
		if cl != nil {
			cl.ParseError = err.Error()
			cl.ParsePartial = true
		}
		return cl, err
	}
	return cl, nil
}

// Convert maps a parsed [syntax.File] into a [cmdast.CommandLine]. f may be nil.
func Convert(f *syntax.File, raw string, lang cmdast.SourceLang) *cmdast.CommandLine {
	cl := &cmdast.CommandLine{Lang: lang, Raw: raw}
	if f == nil || len(f.Stmts) == 0 {
		return cl
	}
	stmts := make([]*cmdast.Statement, 0, len(f.Stmts))
	for _, s := range f.Stmts {
		if st := convertStmt(s, raw); st != nil {
			stmts = append(stmts, st)
		}
	}
	if len(stmts) == 0 {
		return cl
	}
	if len(stmts) == 1 {
		cl.Statements = stmts
		return cl
	}
	cl.Statements = []*cmdast.Statement{foldSemicolonChain(stmts)}
	return cl
}

func foldSemicolonChain(stmts []*cmdast.Statement) *cmdast.Statement {
	acc := stmts[0]
	for i := 1; i < len(stmts); i++ {
		acc = &cmdast.Statement{
			Kind: cmdast.StmtChain,
			Span: cmdast.Span{Start: acc.Span.Start, End: stmts[i].Span.End},
			Chain: &cmdast.Chain{
				Operator: ";",
				Left:     acc,
				Right:    stmts[i],
			},
		}
	}
	return acc
}
