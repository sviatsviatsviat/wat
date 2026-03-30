package shellparse

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// psParseEnvScript is loaded from ps_parse_env.ps1. At runtime [parseComplexOneShot] sets **WAT_CMD** to the
// command string; the script parses it with **System.Management.Automation.Language.Parser** and prints one
// JSON line compatible with [decodePwshJSON] and [core.CommandNode].
//
//go:embed ps_parse_env.ps1
var psParseEnvScript string

// PowerShellParser implements [core.ShellParser] for PowerShell.
//
// It uses a **hybrid** strategy:
//   - **Simple** input (see [detectComplexity]) is handled in-process by [PowerShellParser.parseFast]:
//     pipeline split, tokenizer, and a small known-switch set.
//   - **Complex** input is delegated to **pwsh** via [parseComplexOneShot], which runs the embedded script
//     with the command in the **WAT_CMD** environment variable. If **pwsh** is missing or JSON decoding fails,
//     [PowerShellParser.parseFull] falls back to [PowerShellParser.parseFast] without surfacing an error
//     from the fallback path (the returned [core.ParseResult] still has Dialect [core.DialectPowerShell]).
//
// [PowerShellParser.Parse] always returns a nil error from the fast path; errors from **pwsh** are absorbed by
// falling back when possible.
type PowerShellParser struct{}

// Dialect returns [core.DialectPowerShell].
func (p *PowerShellParser) Dialect() string { return core.DialectPowerShell }

// Parse sets [core.ParseResult.Raw] and [core.ParseResult.Dialect], then fills [core.ParseResult.Pipeline]
// from either the fast path or **pwsh** + JSON decode, with degradation as described on [PowerShellParser].
func (p *PowerShellParser) Parse(raw string) (core.ParseResult, error) {
	result := core.ParseResult{Dialect: core.DialectPowerShell, Raw: raw}
	var nodes []core.CommandNode
	var err error
	if detectComplexity(raw) == complexitySimple {
		nodes, err = p.parseFast(raw)
	} else {
		nodes, err = p.parseFull(raw)
	}
	result.Pipeline = nodes
	return result, err
}

type complexity int

const (
	complexitySimple  complexity = iota // tokenize-only path ([parseFast])
	complexityComplex                   // **pwsh** + embedded parser script ([parseComplexOneShot])
)

// detectComplexity is intentionally conservative: false positives only add a **pwsh** subprocess cost;
// false negatives would mis-parse complex scripts on the fast path.
func detectComplexity(raw string) complexity {
	for _, sub := range [...]string{
		psComplexBraceOpen,
		psComplexSubExpr,
		psComplexArrayLiteral,
		psComplexHashtableLit,
		psComplexHereStringDbl,
		psComplexHereStringSgl,
	} {
		if strings.Contains(raw, sub) {
			return complexityComplex
		}
	}
	if strings.ContainsAny(raw, psComplexParenRunes) && strings.Contains(raw, psComplexPipe) {
		return complexityComplex
	}
	if strings.Contains(raw, psComplexBacktickNewline) || strings.Count(raw, psComplexNewline) > 0 {
		return complexityComplex
	}
	return complexitySimple
}

// parseFull tries [parseComplexOneShot]; on any failure it returns [PowerShellParser.parseFast] instead.
func (p *PowerShellParser) parseFull(raw string) ([]core.CommandNode, error) {
	nodes, err := parseComplexOneShot(raw)
	if err != nil {
		return p.parseFast(raw)
	}
	return nodes, nil
}

// parseComplexOneShot runs **pwsh -NoProfile -NoLogo -Command** with [psParseEnvScript], passing the command
// string in **WAT_CMD** so the shell does not need to re-escape the argument for **-Command**.
func parseComplexOneShot(raw string) ([]core.CommandNode, error) {
	cmd := exec.Command("pwsh", "-NoProfile", "-NoLogo", "-Command", psParseEnvScript)
	cmd.Env = append(os.Environ(), "WAT_CMD="+raw)
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return decodePwshJSON(out)
}

// pwshNodeJSON matches **ConvertTo-Json** property names from the embedded script.
type pwshNodeJSON struct {
	Name       string            `json:"Name"`
	Args       []string          `json:"Args"`
	Flags      map[string]string `json:"Flags"`
	Switches   []string          `json:"Switches"`
	PipeIndex  int               `json:"PipeIndex"`
	PipeLength int               `json:"PipeLength"`
}

// decodePwshJSON unmarshals a JSON array of nodes, or a single object when **ConvertTo-Json** emitted one value.
func decodePwshJSON(out []byte) ([]core.CommandNode, error) {
	payload := strings.TrimSpace(string(out))
	if payload == "" {
		return nil, fmt.Errorf("empty pwsh output")
	}
	rawNodes, err := unmarshalPwshNodeJSONArray(payload)
	if err != nil {
		return nil, err
	}
	nodes := make([]core.CommandNode, len(rawNodes))
	for i, n := range rawNodes {
		if n.Flags == nil {
			n.Flags = map[string]string{}
		}
		nodes[i] = core.CommandNode{
			Name:       n.Name,
			Args:       n.Args,
			Flags:      n.Flags,
			Switches:   n.Switches,
			PipeIndex:  n.PipeIndex,
			PipeLength: n.PipeLength,
		}
	}
	return nodes, nil
}

// unmarshalPwshNodeJSONArray decodes a JSON array of [pwshNodeJSON]. If that fails, it tries a single object
// (PowerShell **ConvertTo-Json** with one element). If both fail, it returns the array unmarshal error.
func unmarshalPwshNodeJSONArray(payload string) ([]pwshNodeJSON, error) {
	var rawNodes []pwshNodeJSON
	errArray := json.Unmarshal([]byte(payload), &rawNodes)
	if errArray == nil {
		return rawNodes, nil
	}
	var one pwshNodeJSON
	errObj := json.Unmarshal([]byte(payload), &one)
	if errObj != nil {
		return nil, errArray
	}
	return []pwshNodeJSON{one}, nil
}

// parseFast splits **raw** on top-level **|** ([splitPipeline]), tokenizes each stage ([tokenizeStage]),
// uses the first token as [core.CommandNode.Name] (no alias resolution), and maps **-Name** tokens into
// [core.CommandNode.Flags] or [core.CommandNode.Switches] using [isKnownSwitch] heuristics.
func (p *PowerShellParser) parseFast(raw string) ([]core.CommandNode, error) {
	stages := splitPipeline(raw)
	nodes := make([]core.CommandNode, 0, len(stages))
	for i, stage := range stages {
		tokens := tokenizeStage(stage)
		if len(tokens) == 0 {
			continue
		}
		node := core.CommandNode{
			Name:       tokens[0],
			Flags:      make(map[string]string),
			PipeIndex:  i,
			PipeLength: len(stages),
		}
		j := 1
		for j < len(tokens) {
			tok := tokens[j]
			if tok == psTokenPipe || tok == psTokenSemicolon {
				j++
				continue
			}
			if strings.HasPrefix(tok, switchPrefix) {
				paramName := tok
				if isKnownSwitch(paramName) {
					node.Switches = append(node.Switches, paramName)
				} else if j+1 < len(tokens) && !strings.HasPrefix(tokens[j+1], switchPrefix) && tokens[j+1] != psTokenPipe && tokens[j+1] != psTokenSemicolon {
					node.Flags[paramName] = tokens[j+1]
					j++
				} else {
					node.Switches = append(node.Switches, paramName)
				}
			} else {
				node.Args = append(node.Args, tok)
			}
			j++
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

// splitPipeline splits **raw** on top-level **|** and **;** (not inside quotes). Quote markers are
// preserved in the stage text for downstream tokenization.
func splitPipeline(raw string) []string {
	var stages []string
	var current strings.Builder
	inSingle, inDouble := false, false

	for _, r := range raw {
		switch {
		case r == psRuneSingleQuote && !inDouble:
			inSingle = !inSingle
			current.WriteRune(r)
		case r == psRuneDoubleQuote && !inSingle:
			inDouble = !inDouble
			current.WriteRune(r)
		case (r == psRunePipe || r == psRuneSemicolon) && !inSingle && !inDouble:
			stages = append(stages, strings.TrimSpace(current.String()))
			current.Reset()
		default:
			current.WriteRune(r)
		}
	}
	if s := strings.TrimSpace(current.String()); s != "" {
		stages = append(stages, s)
	}
	return stages
}

// tokenizeStage splits on whitespace and emits **|** / **;** as separate tokens; supports backtick escape,
// double-quoted segments (with backtick escapes inside), and single-quoted literal segments.
func tokenizeStage(input string) []string {
	var tokens []string
	var current strings.Builder
	runes := []rune(strings.TrimSpace(input))
	i := 0

	flush := func() {
		if current.Len() > 0 {
			tokens = append(tokens, current.String())
			current.Reset()
		}
	}

	for i < len(runes) {
		ch := runes[i]
		switch {
		case ch == psRuneBacktick && i+1 < len(runes):
			i++
			current.WriteRune(runes[i])

		case ch == psRuneDoubleQuote:
			i++
			for i < len(runes) && runes[i] != psRuneDoubleQuote {
				if runes[i] == psRuneBacktick && i+1 < len(runes) {
					i++
				}
				current.WriteRune(runes[i])
				i++
			}

		case ch == psRuneSingleQuote:
			i++
			for i < len(runes) && runes[i] != psRuneSingleQuote {
				current.WriteRune(runes[i])
				i++
			}

		case ch == psRunePipe || ch == psRuneSemicolon:
			flush()
			tokens = append(tokens, string(ch))

		case ch == psRuneSpace || ch == psRuneTab:
			flush()

		default:
			current.WriteRune(ch)
		}
		i++
	}
	flush()
	return tokens
}

// knownSwitches lists parameters treated as boolean switches when the fast path cannot infer a separate value.
var knownSwitches = map[string]bool{
	"-Force": true, "-Recurse": true, "-Verbose": true,
	"-WhatIf": true, "-Confirm": true, "-Debug": true,
	"-NoProfile": true, "-NonInteractive": true,
	"-NoClobber": true, "-Append": true,
}

// isKnownSwitch reports whether **param** is treated as a switch (not a flag with a separate value) on the fast path.
func isKnownSwitch(param string) bool {
	return knownSwitches[param]
}
