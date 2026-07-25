package portio

import (
	"fmt"
	"io"
	"os"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/dialect"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/paths"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig"
)

// Exit codes returned by Run.
const (
	ExitOK             = 0
	ExitRuntimeFailure = 3
)

// Config holds options for Run.
type Config struct {
	// From is the source agent dialect.
	From string
	// To is the target agent dialect.
	To string
	// InputPath is the input config file; empty uses the source agent's default path.
	InputPath string
	// OutputPath is the output file; empty writes to out.
	OutputPath string
}

// Deps holds injectable dependencies for Run.
type Deps struct {
	Getwd     func() (string, error)
	ReadFile  func(string) ([]byte, error)
	WriteFile func(string, []byte, os.FileMode) error
}

// DefaultDeps returns production dependencies backed by the OS.
func DefaultDeps() Deps {
	return Deps{
		Getwd:     os.Getwd,
		ReadFile:  os.ReadFile,
		WriteFile: os.WriteFile,
	}
}

// ParseDialect parses a required --from/--to dialect flag value.
func ParseDialect(value, flag string) (string, error) {
	if value == "" {
		return "", fmt.Errorf("--%s is required (claude, copilot, or cursor)", flag)
	}
	d := dialect.Parse(value)
	if d == "" {
		return "", fmt.Errorf("unknown agent dialect %q for --%s (want claude, copilot, or cursor)", value, flag)
	}
	return d, nil
}

// DefaultInputPath returns the well-known config path for from under wd.
func DefaultInputPath(from, wd string) string {
	return paths.ConfigPath(from, wd)
}

// Run translates a native hook config file from one agent dialect to another.
func Run(cfg Config, deps Deps, out, errOut io.Writer) int {
	inputPath := cfg.InputPath
	if inputPath == "" {
		wd, err := deps.Getwd()
		if err != nil {
			_, _ = fmt.Fprintf(errOut, "wat port: getwd: %v\n", err)
			return ExitRuntimeFailure
		}
		inputPath = DefaultInputPath(cfg.From, wd)
	}

	data, err := deps.ReadFile(inputPath)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat port: read %s: %v\n", inputPath, err)
		return ExitRuntimeFailure
	}

	translated, warns, err := portconfig.Translate(data, cfg.From, cfg.To)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat port: %v\n", err)
		return ExitRuntimeFailure
	}

	for _, w := range warns {
		_, _ = fmt.Fprintf(errOut, "wat port: warning: %s\n", w)
	}

	translated = append(translated, '\n')
	if cfg.OutputPath == "" {
		if _, err := out.Write(translated); err != nil {
			_, _ = fmt.Fprintf(errOut, "wat port: write stdout: %v\n", err)
			return ExitRuntimeFailure
		}
		return ExitOK
	}

	if err := deps.WriteFile(cfg.OutputPath, translated, 0o644); err != nil {
		_, _ = fmt.Fprintf(errOut, "wat port: write %s: %v\n", cfg.OutputPath, err)
		return ExitRuntimeFailure
	}
	return ExitOK
}
