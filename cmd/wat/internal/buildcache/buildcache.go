package buildcache

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// cacheDirName is excluded from the source manifest (build outputs live there).
const cacheDirName = ".cache"

// ManifestArgument requests registration metadata from the generated hook binary.
const ManifestArgument = "__wat_manifest"

const bootstrapSource = `package main

import (
	"encoding/json"
	"fmt"
	"os"

	authored %q
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func main() {
	if len(os.Args) == 2 && os.Args[1] == %q {
		if err := json.NewEncoder(os.Stdout).Encode(run.Inspect(authored.Hooks...)); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "hooks manifest: %%v\n", err)
			os.Exit(1)
		}
		return
	}
	run.Serve(authored.Hooks...)
}
`

// Deps holds injectable filesystem and exec dependencies.
type Deps struct {
	Getenv    func(string) string
	Stat      func(string) (os.FileInfo, error)
	ReadDir   func(string) ([]os.DirEntry, error)
	ReadFile  func(string) ([]byte, error)
	MkdirAll  func(string, os.FileMode) error
	WriteFile func(string, []byte, os.FileMode) error
	Command   func(string, ...string) *exec.Cmd
}

// DefaultDeps returns production dependencies backed by the OS.
func DefaultDeps() Deps {
	return Adapt(os.Getenv, os.Stat, os.ReadDir, os.ReadFile, os.MkdirAll, os.WriteFile, exec.Command)
}

// Adapt builds Deps from the injectable fields shared by CLI packages.
func Adapt(
	getenv func(string) string,
	stat func(string) (os.FileInfo, error),
	readDir func(string) ([]os.DirEntry, error),
	readFile func(string) ([]byte, error),
	mkdirAll func(string, os.FileMode) error,
	writeFile func(string, []byte, os.FileMode) error,
	command func(string, ...string) *exec.Cmd,
) Deps {
	return Deps{
		Getenv:    getenv,
		Stat:      stat,
		ReadDir:   readDir,
		ReadFile:  readFile,
		MkdirAll:  mkdirAll,
		WriteFile: writeFile,
		Command:   command,
	}
}

type buildSettings struct {
	goos       string
	goarch     string
	goflags    string
	cgoEnabled string
	goVersion  string
}

// CacheKey returns the content-addressed cache key for watDir sources and version.
func CacheKey(watDir, version string, deps Deps) (string, error) {
	manifest, err := manifest(watDir, version, deps)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(manifest)
	return hex.EncodeToString(sum[:]), nil
}

func manifest(watDir, version string, deps Deps) ([]byte, error) {
	settings, err := resolveBuildSettings(deps)
	if err != nil {
		return nil, err
	}

	files, err := listSourceFiles(watDir, "", deps)
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	var b bytes.Buffer
	writePart := func(label string, data []byte) {
		_, _ = b.WriteString(label)
		_ = b.WriteByte(0)
		_, _ = b.Write(data)
		_ = b.WriteByte(0)
	}

	writePart("wat_version", []byte(version))
	writePart("goos", []byte(settings.goos))
	writePart("goarch", []byte(settings.goarch))
	writePart("goflags", []byte(settings.goflags))
	writePart("cgo_enabled", []byte(settings.cgoEnabled))
	writePart("go_version", []byte(settings.goVersion))
	writePart("bootstrap", []byte(bootstrapSource))
	writePart("manifest_argument", []byte(ManifestArgument))

	for _, rel := range files {
		data, err := deps.ReadFile(filepath.Join(watDir, rel))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", filepath.Join(watDir, rel), err)
		}
		writePart("file:"+filepath.ToSlash(rel), data)
	}

	return b.Bytes(), nil
}

func resolveBuildSettings(deps Deps) (buildSettings, error) {
	getenv := deps.Getenv
	if getenv == nil {
		getenv = os.Getenv
	}
	goos := strings.TrimSpace(getenv("GOOS"))
	if goos == "" {
		goos = runtime.GOOS
	}
	goarch := strings.TrimSpace(getenv("GOARCH"))
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	goVersion, err := goEnv(deps, "GOVERSION")
	if err != nil {
		return buildSettings{}, err
	}
	return buildSettings{
		goos:       goos,
		goarch:     goarch,
		goflags:    getenv("GOFLAGS"),
		cgoEnabled: getenv("CGO_ENABLED"),
		goVersion:  goVersion,
	}, nil
}

func goEnv(deps Deps, key string) (string, error) {
	cmd := deps.Command("go", "env", key)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("go env %s: %w", key, err)
	}
	return strings.TrimSpace(string(out)), nil
}

func listSourceFiles(dir, rel string, deps Deps) ([]string, error) {
	entries, err := deps.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}
	var files []string
	for _, e := range entries {
		name := e.Name()
		childRel := name
		if rel != "" {
			childRel = filepath.Join(rel, name)
		}
		if e.IsDir() {
			if name == cacheDirName {
				continue
			}
			sub, err := listSourceFiles(filepath.Join(dir, name), childRel, deps)
			if err != nil {
				return nil, err
			}
			files = append(files, sub...)
			continue
		}
		files = append(files, childRel)
	}
	return files, nil
}

// BinaryPath returns the content-addressed hooks binary path for key.
func BinaryPath(watDir, key string) string {
	name := "hooks"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(watDir, cacheDirName, key, name)
}

// ErrBuildFailed indicates go build failed while ensuring the hooks binary.
var ErrBuildFailed = errors.New("hook binary build failed")

// Ensure returns the content-addressed hooks binary path, building it on cache miss.
// Build diagnostics are written to errOut. On failure, err is ErrBuildFailed or a
// wrapped runtime error.
func Ensure(watDir, version string, deps Deps, errOut io.Writer) (binPath string, err error) {
	key, err := CacheKey(watDir, version, deps)
	if err != nil {
		return "", err
	}
	binPath = BinaryPath(watDir, key)

	if _, err := deps.Stat(binPath); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("stat %s: %w", binPath, err)
		}
		if !build(watDir, binPath, deps, errOut) {
			return "", ErrBuildFailed
		}
	}
	return binPath, nil
}

func build(watDir, binPath string, deps Deps, errOut io.Writer) bool {
	cacheDir := filepath.Dir(binPath)
	if err := deps.MkdirAll(cacheDir, 0o755); err != nil {
		_, _ = fmt.Fprintf(errOut, "wat run: create %s: %v\n", cacheDir, err)
		return false
	}

	settings, err := resolveBuildSettings(deps)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat run: %v\n", err)
		return false
	}

	modulePath, err := currentModulePath(watDir, deps, settings)
	if err != nil {
		_, _ = fmt.Fprintf(errOut, "wat run: %v\n", err)
		return false
	}

	bootstrapDir := filepath.Join(cacheDir, "bootstrap")
	if err := deps.MkdirAll(bootstrapDir, 0o755); err != nil {
		_, _ = fmt.Fprintf(errOut, "wat run: create %s: %v\n", bootstrapDir, err)
		return false
	}
	source := fmt.Sprintf(bootstrapSource, modulePath, ManifestArgument)
	writeFile := deps.WriteFile
	if writeFile == nil {
		writeFile = os.WriteFile
	}
	if err := writeFile(filepath.Join(bootstrapDir, "main.go"), []byte(source), 0o600); err != nil {
		_, _ = fmt.Fprintf(errOut, "wat run: write bootstrap: %v\n", err)
		return false
	}

	cmd := deps.Command("go", "build", "-o", binPath, bootstrapDir)
	cmd.Dir = watDir
	cmd.Env = pinnedBuildEnv(os.Environ(), settings)
	out, err := cmd.CombinedOutput()
	if err == nil {
		return true
	}

	if len(out) > 0 {
		_, _ = errOut.Write(out)
		if out[len(out)-1] != '\n' {
			_, _ = fmt.Fprintln(errOut, "")
		}
	}
	_, _ = fmt.Fprintf(errOut, "wat run: go build failed in %s\n", watDir)
	return false
}

func currentModulePath(watDir string, deps Deps, settings buildSettings) (string, error) {
	cmd := deps.Command("go", "list", "-m", "-f={{.Path}}")
	cmd.Dir = watDir
	cmd.Env = pinnedBuildEnv(os.Environ(), settings)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("resolve hook module path: %w", err)
	}
	modulePath := strings.TrimSpace(string(out))
	if modulePath == "" {
		return "", fmt.Errorf("resolve hook module path: empty module path")
	}
	return modulePath, nil
}

func pinnedBuildEnv(base []string, settings buildSettings) []string {
	env := append([]string(nil), base...)
	env = setEnv(env, "GOOS", settings.goos)
	env = setEnv(env, "GOARCH", settings.goarch)
	env = setEnv(env, "GOFLAGS", settings.goflags)
	env = setEnv(env, "CGO_ENABLED", settings.cgoEnabled)
	return env
}

func setEnv(env []string, key, value string) []string {
	prefix := key + "="
	for i, e := range env {
		if strings.HasPrefix(e, prefix) {
			env[i] = prefix + value
			return env
		}
	}
	return append(env, prefix+value)
}
