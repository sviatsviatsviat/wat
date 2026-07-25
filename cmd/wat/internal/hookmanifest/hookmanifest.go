package hookmanifest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/buildcache"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// LoadTimeout bounds how long LoadBinary waits for manifest mode to finish.
var LoadTimeout = 30 * time.Second

// ErrLoadTimeout indicates the authored hook binary did not finish manifest mode in time.
var ErrLoadTimeout = errors.New("hook manifest timed out")

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
	if err := runWithTimeout(cmd, LoadTimeout); err != nil {
		if errors.Is(err, ErrLoadTimeout) {
			return run.Manifest{}, fmt.Errorf("read hook manifest: %w after %s", ErrLoadTimeout, LoadTimeout)
		}
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

func runWithTimeout(cmd *exec.Cmd, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := cmd.Start(); err != nil {
		return err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		<-done
		return ErrLoadTimeout
	}
}
