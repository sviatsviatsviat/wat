package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

type runDeps struct {
	getenv   func(string) string
	getwd    func() (string, error)
	stat     func(string) (os.FileInfo, error)
	readDir  func(string) ([]os.DirEntry, error)
	readFile func(string) ([]byte, error)
	mkdirAll func(string, os.FileMode) error
	command  func(string, ...string) *exec.Cmd
	runCmd   func(*exec.Cmd) error
}

func defaultRunDeps() runDeps {
	return runDeps{
		getenv:   os.Getenv,
		getwd:    os.Getwd,
		stat:     os.Stat,
		readDir:  os.ReadDir,
		readFile: os.ReadFile,
		mkdirAll: os.MkdirAll,
		command:  exec.Command,
		runCmd: func(cmd *exec.Cmd) error {
			return cmd.Run()
		},
	}
}

type runConfig struct {
	agent      string
	event      string
	failClosed bool
}

func runHook(cfg runConfig, deps runDeps) int {
	watDir, err := resolveWatDir(deps)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "wat run: %v\n", err)
		return exitRuntimeFailure
	}

	key, err := hookBuildCacheKey(watDir, deps)
	if err != nil {
		_, _ = fmt.Fprintf(stderr, "wat run: %v\n", err)
		return exitRuntimeFailure
	}
	binPath := hooksBinaryPath(watDir, key)

	if _, err := deps.stat(binPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			_, _ = fmt.Fprintf(stderr, "wat run: stat %s: %v\n", binPath, err)
			return exitRuntimeFailure
		}
		if ok := buildHookBinary(watDir, binPath, deps); !ok {
			if cfg.failClosed {
				return exitFailClosed
			}
			return exitBuildFailed
		}
	}

	cmd := deps.command(binPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if len(cmd.Env) == 0 {
		cmd.Env = append([]string(nil), os.Environ()...)
	}
	if cfg.agent != "" {
		cmd.Env = append(cmd.Env, "WAT_AGENT="+cfg.agent)
	}
	if cfg.event != "" {
		cmd.Env = append(cmd.Env, "WAT_EVENT="+cfg.event)
	}

	if err := deps.runCmd(cmd); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		_, _ = fmt.Fprintf(stderr, "wat run: exec %s: %v\n", binPath, err)
		return exitRuntimeFailure
	}
	return exitOK
}

func hookBuildCacheKey(watDir string, deps runDeps) (string, error) {
	manifest, err := hookBuildManifest(watDir, deps)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(manifest)
	return hex.EncodeToString(sum[:]), nil
}

func hookBuildManifest(watDir string, deps runDeps) ([]byte, error) {
	entries, err := deps.readDir(watDir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", watDir, err)
	}

	var goFiles []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".go") {
			goFiles = append(goFiles, name)
		}
	}
	sort.Strings(goFiles)

	var b bytes.Buffer
	writePart := func(label string, data []byte) {
		_, _ = b.WriteString(label)
		_ = b.WriteByte(0)
		_, _ = b.Write(data)
		_ = b.WriteByte(0)
	}

	writePart("wat_version", []byte(watModuleVersionFn()))
	writePart("goos", []byte(runtime.GOOS))
	writePart("goarch", []byte(runtime.GOARCH))

	goMod, err := deps.readFile(filepath.Join(watDir, "go.mod"))
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", filepath.Join(watDir, "go.mod"), err)
	}
	writePart("go.mod", goMod)

	goSum, err := deps.readFile(filepath.Join(watDir, "go.sum"))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("read %s: %w", filepath.Join(watDir, "go.sum"), err)
		}
		goSum = nil
	}
	writePart("go.sum", goSum)

	for _, name := range goFiles {
		data, err := deps.readFile(filepath.Join(watDir, name))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", filepath.Join(watDir, name), err)
		}
		writePart("go:"+name, data)
	}

	return b.Bytes(), nil
}

func resolveWatDir(deps runDeps) (string, error) {
	root := strings.TrimSpace(deps.getenv("WAT_PROJECT_DIR"))
	if root != "" {
		abs, err := filepath.Abs(root)
		if err != nil {
			return "", fmt.Errorf("resolve WAT_PROJECT_DIR: %w", err)
		}
		return watDirFromRoot(abs, deps)
	}
	cwd, err := deps.getwd()
	if err != nil {
		return "", fmt.Errorf("getwd: %w", err)
	}
	return findWatDir(cwd, deps)
}

func watDirFromRoot(root string, deps runDeps) (string, error) {
	watDir := filepath.Join(root, ".wat")
	if err := mustHaveWatFiles(watDir, deps); err != nil {
		return "", fmt.Errorf("%s is not a wat hook project: %w", watDir, err)
	}
	return watDir, nil
}

func findWatDir(start string, deps runDeps) (string, error) {
	dir := start
	for {
		watDir := filepath.Join(dir, ".wat")
		if err := mustHaveWatFiles(watDir, deps); err == nil {
			return watDir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("no .wat/ project found from %s (run \"wat init\" first)", start)
		}
		dir = parent
	}
}

func mustHaveWatFiles(watDir string, deps runDeps) error {
	hooksPath := filepath.Join(watDir, "hooks.go")
	if _, err := deps.stat(hooksPath); err != nil {
		return fmt.Errorf("stat %s: %w", hooksPath, err)
	}
	goModPath := filepath.Join(watDir, "go.mod")
	if _, err := deps.stat(goModPath); err != nil {
		return fmt.Errorf("stat %s: %w", goModPath, err)
	}
	return nil
}

func hooksBinaryPath(watDir, key string) string {
	name := "hooks"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(watDir, ".cache", key, name)
}

func buildHookBinary(watDir, binPath string, deps runDeps) bool {
	cacheDir := filepath.Dir(binPath)
	if err := deps.mkdirAll(cacheDir, 0o755); err != nil {
		_, _ = fmt.Fprintf(stderr, "wat run: create %s: %v\n", cacheDir, err)
		return false
	}

	cmd := deps.command("go", "build", "-o", binPath)
	cmd.Dir = watDir
	out, err := cmd.CombinedOutput()
	if err == nil {
		return true
	}

	if len(out) > 0 {
		_, _ = stderr.Write(out)
		if out[len(out)-1] != '\n' {
			_, _ = fmt.Fprintln(stderr, "")
		}
	}
	_, _ = fmt.Fprintf(stderr, "wat run: go build failed in %s\n", watDir)
	return false
}

