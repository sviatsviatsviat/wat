package execcommand

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// commandForArgs returns an [exec.Cmd] for args. It returns an error if args is empty.
// If args[0] resolves with [exec.LookPath], the process is started directly; otherwise a shell
// runs the same arguments (shell builtins and non-PATH commands like Windows "echo").
func commandForArgs(args []string) (*exec.Cmd, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("no command arguments provided")
	}
	if _, err := exec.LookPath(args[0]); err == nil {
		return exec.Command(args[0], args[1:]...), nil
	}
	return shellCommand(args)
}

func shellCommand(args []string) (*exec.Cmd, error) {
	if runtime.GOOS == "windows" {
		cmdExe := comspecOrCmd()
		escaped := make([]string, len(args))
		for i, a := range args {
			escaped[i] = windowsShellArg(a)
		}
		cmdArgs := append([]string{"/C"}, escaped...)
		return exec.Command(cmdExe, cmdArgs...), nil
	}
	shPath, err := exec.LookPath("sh")
	if err != nil {
		return nil, err
	}
	return exec.Command(shPath, "-c", posixShellLine(args)), nil
}

// windowsShellArg returns a string safe for one token after cmd.exe /C. Whether the argument
// needs double quotes is decided from the original string (whitespace or embedded "). Inside
// quotes, metacharacters &|<> are literal; " is doubled, % is doubled (%%), and ! and ^ are
// caret-escaped for delayed expansion / caret semantics. Without quotes, ^ and cmd metacharacters
// are caret-escaped, including parentheses.
func windowsShellArg(s string) string {
	if s == "" {
		return `""`
	}

	needsQuotes := strings.Contains(s, `"`)
	if !needsQuotes {
		for _, r := range s {
			if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
				needsQuotes = true
				break
			}
		}
	}

	var b strings.Builder
	b.Grow(len(s) + 8)

	if needsQuotes {
		for _, r := range s {
			switch r {
			case '"':
				b.WriteString(`""`)
			case '%':
				b.WriteString("%%")
			case '!':
				b.WriteString("^!")
			case '^':
				b.WriteString("^^")
			default:
				b.WriteRune(r)
			}
		}
		return `"` + b.String() + `"`
	}

	for _, r := range s {
		switch r {
		case '^':
			b.WriteString("^^")
		case '&', '|', '<', '>', '%', '(', ')':
			b.WriteRune('^')
			b.WriteRune(r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

func comspecOrCmd() string {
	if c := os.Getenv("COMSPEC"); c != "" {
		return c
	}
	if p, err := exec.LookPath("cmd.exe"); err == nil {
		return p
	}
	return "cmd.exe"
}

// posixShellLine joins arguments into a string safe for sh -c using single-quote escaping.
func posixShellLine(args []string) string {
	quoted := make([]string, len(args))
	for i, a := range args {
		quoted[i] = posixSingleQuoted(a)
	}
	return strings.Join(quoted, " ")
}

func posixSingleQuoted(s string) string {
	// End quote, escaped ', resume quote — portable for sh -c.
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}
