package doctor

import (
	"io"
	"os"
	"os/exec"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/buildcache"
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/hookmanifest"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Deps holds injectable dependencies for doctor checks.
type Deps struct {
	Getenv       func(string) string
	LookupEnv    func(string) (string, bool)
	Getwd        func() (string, error)
	Stat         func(string) (os.FileInfo, error)
	ReadDir      func(string) ([]os.DirEntry, error)
	ReadFile     func(string) ([]byte, error)
	MkdirAll     func(string, os.FileMode) error
	WriteFile    func(string, []byte, os.FileMode) error
	Remove       func(string) error
	LookPath     func(string) (string, error)
	Command      func(string, ...string) *exec.Cmd
	IsTerminal   func(io.Writer) bool
	LoadManifest func(string, string, buildcache.Deps, io.Writer) (run.Manifest, error)

	// WatVersion is included in the hook build cache key (typically the wat module version).
	WatVersion string
}

// DefaultDeps returns production dependencies backed by the OS.
func DefaultDeps() Deps {
	return Deps{
		Getenv:    os.Getenv,
		LookupEnv: os.LookupEnv,
		Getwd:     os.Getwd,
		Stat:      os.Stat,
		ReadDir:   os.ReadDir,
		ReadFile:  os.ReadFile,
		MkdirAll:  os.MkdirAll,
		WriteFile: os.WriteFile,
		Remove:    os.Remove,
		LookPath:  exec.LookPath,
		Command:   exec.Command,
		IsTerminal: func(w io.Writer) bool {
			f, ok := w.(*os.File)
			if !ok {
				return false
			}
			info, err := f.Stat()
			return err == nil && info.Mode()&os.ModeCharDevice != 0
		},
		LoadManifest: hookmanifest.Load,
	}
}
