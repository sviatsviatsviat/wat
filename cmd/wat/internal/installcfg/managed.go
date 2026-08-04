package installcfg

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"
)

// ParseWatRunFlags extracts --agent and --event from a wat run shell command.
func ParseWatRunFlags(command string) (agent, event string, ok bool) {
	fields := splitCommandFields(command)
	for i := 0; i < len(fields)-1; i++ {
		switch fields[i] {
		case "--agent":
			agent = fields[i+1]
		case "--event":
			event = fields[i+1]
		}
	}
	return agent, event, agent != "" && event != ""
}

// IsWatManagedCommand reports whether command is a wat-managed hook entry.
func IsWatManagedCommand(command, agent, event, watAbs string) bool {
	parsedAgent, parsedEvent, ok := ParseWatRunFlags(command)
	if !ok || parsedAgent != agent || parsedEvent != event {
		return false
	}
	return isWatRunExecutable(command, watAbs)
}

// IsWatManagedAgentCommand reports whether command is a wat-managed hook entry
// for agent, regardless of its native event.
func IsWatManagedAgentCommand(command, agent, watAbs string) bool {
	parsedAgent, _, ok := ParseWatRunFlags(command)
	if !ok || parsedAgent != agent {
		return false
	}
	return isWatRunExecutable(command, watAbs)
}

func isWatRunExecutable(command, watAbs string) bool {
	fields := splitCommandFields(strings.TrimSpace(command))
	if len(fields) != 6 || fields[1] != "run" {
		return false
	}
	program := fields[0]
	watAbs = strings.TrimSpace(watAbs)
	if program == watAbs {
		return true
	}
	base := strings.TrimSuffix(strings.ToLower(commandBaseName(program)), ".exe")
	return base == "wat"
}

// WatRunCommand builds the shell command wat install writes for an agent event.
// Paths that contain whitespace or shell metacharacters are double-quoted with
// escapes so host shells invoke the installed binary literally.
func WatRunCommand(watAbs, agent, event string) string {
	return fmt.Sprintf("%s run --agent %s --event %s", quoteShellArg(watAbs), agent, event)
}

func quoteShellArg(s string) string {
	if s == "" {
		return `""`
	}
	if !needsShellQuotes(s) {
		return s
	}
	var b strings.Builder
	b.Grow(len(s) + 2)
	b.WriteByte('"')
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '\\', '"', '$', '`':
			b.WriteByte('\\')
		}
		b.WriteByte(s[i])
	}
	b.WriteByte('"')
	return b.String()
}

func needsShellQuotes(s string) bool {
	for _, r := range s {
		if unicode.IsSpace(r) {
			return true
		}
		switch r {
		case '"', '\'', '$', '`', '\\', ';', '|', '&', '(', ')', '<', '>', '*', '?', '[', ']', '{', '}', '~', '#', '!':
			return true
		}
	}
	return false
}

// commandBaseName returns the final path element, treating both / and \ as
// separators so Windows install commands parse correctly on any GOOS.
func commandBaseName(program string) string {
	program = strings.ReplaceAll(program, "\\", "/")
	return filepath.Base(program)
}

// splitCommandFields splits a shell-like command into arguments, honoring
// simple single- and double-quoted segments used by wat install.
func splitCommandFields(command string) []string {
	var fields []string
	var b strings.Builder
	inQuote := false
	var quote byte

	flush := func() {
		if b.Len() == 0 {
			return
		}
		fields = append(fields, b.String())
		b.Reset()
	}

	for i := 0; i < len(command); i++ {
		c := command[i]
		if inQuote {
			if c == quote {
				inQuote = false
				continue
			}
			if c == '\\' && quote == '"' && i+1 < len(command) {
				next := command[i+1]
				switch next {
				case '"', '\\', '$', '`':
					b.WriteByte(next)
					i++
					continue
				}
			}
			b.WriteByte(c)
			continue
		}
		switch c {
		case ' ', '\t', '\n', '\r':
			flush()
		case '"', '\'':
			inQuote = true
			quote = c
		default:
			b.WriteByte(c)
		}
	}
	flush()
	return fields
}
