package execcommand

import (
	"path/filepath"
	"regexp"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	"github.com/sviatsviatsviat/wat/internal/cursor"
)

// execHookHandlerAfterFileEdit runs exec templates for afterFileEdit and afterTabFileEdit hooks.
// Optional -f/--file-pattern applies only here: it matches the event file_path from the hook payload.
type execHookHandlerAfterFileEdit struct {
	execHookHandlerBase
	filePathFilterRegexp *regexp.Regexp
	hook                 *cursor.AfterFileEditCursorHookAdapter
}

func (h execHookHandlerAfterFileEdit) Handle() core.HookHandlerResult {
	if h.filePathFilterRegexp != nil {
		if event := h.hook.EventSpecificInput; event != nil {
			normalizedFilePath := filepath.ToSlash(filepath.Clean(event.FilePath))
			if !h.filePathFilterRegexp.MatchString(normalizedFilePath) {
				h.hook.ReturnEmpty()
				return core.HookHandlerResult{Code: cli.ExitSuccess}
			}
		}
	}

	bindings := templateBindingsFromCursorEventPayload(h.hook.CommonInput, h.hook.EventSpecificInput, afterFileEditPlaceholderExtractors)
	return h.runExecWithBindings(bindings, h.hook)
}
