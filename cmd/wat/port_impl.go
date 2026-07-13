package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sviatsviatsviat/wat/agenthooks"
	"github.com/sviatsviatsviat/wat/agenthooks/portconfig"
)

type portConfig struct {
	from       agenthooks.Dialect
	to         agenthooks.Dialect
	inputPath  string
	outputPath string
}

type portDeps struct {
	getwd     func() (string, error)
	readFile  func(string) ([]byte, error)
	writeFile func(string, []byte, os.FileMode) error
}

func defaultPortDeps() portDeps {
	return portDeps{
		getwd:     os.Getwd,
		readFile:  os.ReadFile,
		writeFile: os.WriteFile,
	}
}

func defaultInputPath(from agenthooks.Dialect, wd string) string {
	switch from {
	case agenthooks.Claude:
		return filepath.Join(wd, ".claude", "settings.json")
	case agenthooks.Copilot:
		return filepath.Join(wd, ".github", "hooks", "wat.json")
	case agenthooks.Cursor:
		return filepath.Join(wd, ".cursor", "hooks.json")
	default:
		return ""
	}
}

func parsePortDialect(value, flag string) (agenthooks.Dialect, error) {
	if value == "" {
		return agenthooks.Unknown, fmt.Errorf("--%s is required (claude, copilot, or cursor)", flag)
	}
	d := agenthooks.ParseDialect(value)
	if d == agenthooks.Unknown {
		return agenthooks.Unknown, fmt.Errorf("unknown agent dialect %q for --%s (want claude, copilot, or cursor)", value, flag)
	}
	return d, nil
}

func portProject(cfg portConfig, deps portDeps) int {
	inputPath := cfg.inputPath
	if inputPath == "" {
		wd, err := deps.getwd()
		if err != nil {
			_, _ = fmt.Fprintf(stderr, "wat port: getwd: %v\n", err)
			return exitRuntimeFailure
		}
		inputPath = defaultInputPath(cfg.from, wd)
	}

	data, err := deps.readFile(inputPath)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "wat port: read %s: %v\n", inputPath, err)
		return exitRuntimeFailure
	}

	out, warns, err := portconfig.Translate(data, cfg.from, cfg.to)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "wat port: %v\n", err)
		return exitRuntimeFailure
	}

	for _, w := range warns {
		_, _ = fmt.Fprintf(stderr, "wat port: warning: %s\n", w)
	}

	out = append(out, '\n')
	if cfg.outputPath == "" {
		_, _ = stdout.Write(out)
		return exitOK
	}

	if err := deps.writeFile(cfg.outputPath, out, 0o644); err != nil {
		_, _ = fmt.Fprintf(stderr, "wat port: write %s: %v\n", cfg.outputPath, err)
		return exitRuntimeFailure
	}
	return exitOK
}
