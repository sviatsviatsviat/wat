package initproj

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoMod_requiresVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		wantErr string
		want    []string
	}{
		{
			name:    "empty",
			version: "",
			wantErr: "wat module version",
		},
		{
			name:    "valid",
			version: "v0.0.0-test-000000000000",
			want: []string{
				"github.com/sviatsviatsviat/wat v0.0.0-test-000000000000",
				"module wat-hooks",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GoMod(tt.version)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error for empty version")
				}
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("error = %q", err.Error())
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Fatalf("go.mod = %q, missing %q", got, want)
				}
			}
		})
	}
}

func TestInit_createsScaffold(t *testing.T) {
	dir := t.TempDir()
	var gotCmdName string
	var gotCmdArgs []string
	var gotCmd *exec.Cmd
	deps := DefaultDeps()
	deps.Command = func(name string, args ...string) *exec.Cmd {
		gotCmdName = name
		gotCmdArgs = append([]string(nil), args...)
		gotCmd = exec.Command("go", "version")
		return gotCmd
	}

	if err := Init(dir, false, "v0.0.0-test-000000000000", deps, io.Discard, io.Discard); err != nil {
		t.Fatalf("Init: %v", err)
	}

	hooks, err := os.ReadFile(filepath.Join(dir, ".wat", "hooks.go"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(hooks), "agnostic.UseHooks()") {
		t.Fatalf("hooks.go missing UseHooks: %s", hooks)
	}
	goMod, err := os.ReadFile(filepath.Join(dir, ".wat", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(goMod), "github.com/sviatsviatsviat/wat v0.0.0-test-000000000000") {
		t.Fatalf("go.mod = %s", goMod)
	}
	if gotCmdName != "go" || len(gotCmdArgs) != 2 || gotCmdArgs[0] != "mod" || gotCmdArgs[1] != "tidy" {
		t.Fatalf("expected go mod tidy, got %s %v", gotCmdName, gotCmdArgs)
	}
	if gotCmd == nil || gotCmd.Dir != filepath.Join(dir, ".wat") {
		t.Fatalf("tidy Dir = %v", gotCmd)
	}
}

func TestInit_refusesOverwriteWithoutForce(t *testing.T) {
	dir := t.TempDir()
	deps := DefaultDeps()
	deps.Command = func(string, ...string) *exec.Cmd {
		return exec.Command("go", "version")
	}

	if err := Init(dir, false, "v0.0.0-test", deps, io.Discard, io.Discard); err != nil {
		t.Fatalf("first Init: %v", err)
	}
	err := Init(dir, false, "v0.0.0-test", deps, io.Discard, io.Discard)
	if err == nil {
		t.Fatal("expected overwrite error")
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Fatalf("error = %q", err.Error())
	}
	if err := Init(dir, true, "v0.0.0-test", deps, io.Discard, io.Discard); err != nil {
		t.Fatalf("Init with force: %v", err)
	}
}
