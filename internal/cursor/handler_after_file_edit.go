package cursor

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/sviatsviatsviat/wat/internal/cli"
	"github.com/sviatsviatsviat/wat/internal/core"
	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
)

type afterFileEditHookHandler struct {
	inner       *cursorcore.EventHookHandler
	fields      hookDataAfterFileEditFields
	filePattern *regexp.Regexp
}

func newAfterFileEditHookHandler(rawJSON []byte, commonData cursorcore.HookDataCommon, execCtx core.WatExecutionContext) (core.HookHandler, error) {
	hookData, err := cursorcore.NewHookDataWithCommon[hookDataAfterFileEditFields](rawJSON, commonData)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor afterFileEdit payload: %w", err)
	}
	inner := cursorcore.NewEventHookHandlerFromExtractedFields(
		hookData.HookDataCommon,
		hookData.Fields,
		afterFileEditPlaceholderExtractors,
	)
	var fileRe *regexp.Regexp
	if ptr := execCtx.FilePattern(); ptr != nil {
		re, compileErr := regexp.Compile(*ptr)
		if compileErr != nil {
			return nil, fmt.Errorf("invalid --file-pattern regexp: %w", compileErr)
		}
		fileRe = re
	}
	return afterFileEditHookHandler{
		inner:       inner,
		fields:      hookData.Fields,
		filePattern: fileRe,
	}, nil
}

func (handler afterFileEditHookHandler) Handle(cmd core.Command) core.HookHandlerResult {
	if handler.filePattern != nil {
		normalized := filepath.ToSlash(filepath.Clean(handler.fields.FilePath))
		if !handler.filePattern.MatchString(normalized) {
			return core.HookHandlerResult{Code: cli.ExitSuccess, Output: cursorcore.DefaultHookResponseLine}
		}
	}
	return handler.inner.Handle(cmd)
}
