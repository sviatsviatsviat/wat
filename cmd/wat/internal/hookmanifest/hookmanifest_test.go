package hookmanifest

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"

	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestLoadBinary_decodesManifest(t *testing.T) {
	if os.Getenv("GO_WANT_MANIFEST_HELPER") == "1" {
		_, _ = fmt.Fprint(os.Stdout, `{"version":1,"registrations":[{"dialect":"cursor","event":"beforeShellExecution","handler_count":2}]}`)
		os.Exit(0)
	}
	command := func(string, ...string) *exec.Cmd {
		cmd := exec.Command(os.Args[0], "-test.run=TestLoadBinary_decodesManifest")
		cmd.Env = append(os.Environ(), "GO_WANT_MANIFEST_HELPER=1")
		return cmd
	}

	got, err := LoadBinary("hooks", command)
	if err != nil {
		t.Fatal(err)
	}
	want := run.Manifest{
		Version: 1,
		Registrations: []run.Registration{{
			Dialect:      "cursor",
			Event:        "beforeShellExecution",
			HandlerCount: 2,
		}},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("LoadBinary() = %#v, want %#v", got, want)
	}
}

func TestLoadBinary_rejectsInvalidRegistration(t *testing.T) {
	if os.Getenv("GO_WANT_INVALID_MANIFEST_HELPER") == "1" {
		_, _ = fmt.Fprint(os.Stdout, `{"version":1,"registrations":[{"dialect":"cursor","event":"","handler_count":0}]}`)
		os.Exit(0)
	}
	command := func(string, ...string) *exec.Cmd {
		cmd := exec.Command(os.Args[0], "-test.run=TestLoadBinary_rejectsInvalidRegistration")
		cmd.Env = append(os.Environ(), "GO_WANT_INVALID_MANIFEST_HELPER=1")
		return cmd
	}

	if _, err := LoadBinary("hooks", command); err == nil {
		t.Fatal("LoadBinary() error = nil, want invalid registration error")
	}
}

func TestLoadBinary_timesOut(t *testing.T) {
	command := func(string, ...string) *exec.Cmd {
		return exec.Command("sleep", "10")
	}

	old := LoadTimeout
	LoadTimeout = 50 * time.Millisecond
	defer func() { LoadTimeout = old }()

	_, err := LoadBinary("hooks", command)
	if !errors.Is(err, ErrLoadTimeout) {
		t.Fatalf("LoadBinary() error = %v, want %v", err, ErrLoadTimeout)
	}
}
