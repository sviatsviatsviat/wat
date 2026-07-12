package execcommand

import (
	"path/filepath"
	"regexp"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
)

// execHookHandlerAfterFileEdit runs exec templates for afterFileEdit and afterTabFileEdit hooks.
// Optional -f/--file-pattern applies only here: it matches the event file_path from the hook payload.
type execHookHandlerAfterFileEdit struct {
	execHookHandlerBase
	filePathFilterRegexp *regexp.Regexp
	hook                 core.AfterFileEditHook
}

func (h execHookHandlerAfterFileEdit) Handle() core.HookHandlerResult {
	if h.filePathFilterRegexp != nil {
		normalizedFilePath := filepath.ToSlash(filepath.Clean(h.hook.FilePath()))
		if !h.filePathFilterRegexp.MatchString(normalizedFilePath) {
			h.hook.WriteDefaultToHost()
			return core.HookHandlerResult{Code: cli.ExitSuccess}
		}
	}

	bindings := templateBindingsAfterFileEdit{hook: h.hook}
	return h.runExecWithBindings(bindings, h.hook)
}
