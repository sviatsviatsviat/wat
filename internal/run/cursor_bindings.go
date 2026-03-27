package run

import (
	"fmt"

	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
	"github.com/sviatsviatsviat/wat/internal/template"
)

func templateBindingsForCursor(parsed any) (template.TemplateBindings, error) {
	if parsed == nil {
		return nil, fmt.Errorf("cursor hook run data is nil")
	}
	switch data := parsed.(type) {
	case *cursorcore.CursorHookRunData[cursorcore.AfterFileEditFields]:
		return templateBindingsFromCursorEventPayload(data, afterFileEditPlaceholderExtractors,
			"cursor afterFileEdit payload missing event fields")
	case *cursorcore.CursorHookRunData[cursorcore.AfterShellExecutionFields]:
		return templateBindingsFromCursorEventPayload(data, afterShellExecutionPlaceholderExtractors,
			"cursor afterShellExecution payload missing event fields")
	case *cursorcore.CursorHookRunData[struct{}]:
		return newTemplateBindingsCommon(data.Common), nil
	default:
		return nil, fmt.Errorf("unexpected cursor ParsedData type %T for run", parsed)
	}
}
