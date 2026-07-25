package hookmanifest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/buildcache"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Load ensures the authored hook binary and returns its registration manifest.
func Load(watDir, version string, deps buildcache.Deps, errOut io.Writer) (run.Manifest, error) {
	_, manifest, err := EnsureAndLoad(watDir, version, deps, errOut)
	return manifest, err
}

// EnsureAndLoad ensures the authored hook binary and returns its path and manifest.
func EnsureAndLoad(watDir, version string, deps buildcache.Deps, errOut io.Writer) (string, run.Manifest, error) {
	binPath, err := buildcache.Ensure(watDir, version, deps, errOut)
	if err != nil {
		return "", run.Manifest{}, err
	}
	manifest, err := LoadBinary(binPath, deps.Command)
	if err != nil {
		return "", run.Manifest{}, err
	}
	return binPath, manifest, nil
}

// LoadBinary executes an authored hook binary in manifest mode.
func LoadBinary(binPath string, command func(string, ...string) *exec.Cmd) (run.Manifest, error) {
	cmd := command(binPath, buildcache.ManifestArgument)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return run.Manifest{}, fmt.Errorf("read hook manifest: %w: %s", err, bytes.TrimSpace(stderr.Bytes()))
		}
		return run.Manifest{}, fmt.Errorf("read hook manifest: %w", err)
	}

	var manifest run.Manifest
	if err := json.Unmarshal(stdout.Bytes(), &manifest); err != nil {
		return run.Manifest{}, fmt.Errorf("decode hook manifest: %w", err)
	}
	if manifest.Version != run.ManifestVersion {
		return run.Manifest{}, fmt.Errorf("decode hook manifest: unsupported version %d", manifest.Version)
	}
	for _, registration := range manifest.Registrations {
		if registration.Dialect == "" || registration.Event == "" || registration.HandlerCount < 1 {
			return run.Manifest{}, fmt.Errorf("decode hook manifest: invalid registration")
		}
	}
	return manifest, nil
}
