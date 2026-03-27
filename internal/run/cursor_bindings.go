package run

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/cursor"
)

func templateBindingsForCursor(parsed any) (templateBindings, error) {
	if parsed == nil {
		return nil, fmt.Errorf("cursor hook run data is nil")
	}
	switch data := parsed.(type) {
	case *cursor.CursorHookRunData[cursor.AfterFileEditFields]:
		return templateBindingsFromCursorEventPayload(data, afterFileEditPlaceholderExtractors,
			"cursor afterFileEdit payload missing event fields")
	case *cursor.CursorHookRunData[cursor.AfterShellExecutionFields]:
		return templateBindingsFromCursorEventPayload(data, afterShellExecutionPlaceholderExtractors,
			"cursor afterShellExecution payload missing event fields")
	case *cursor.CursorHookRunData[struct{}]:
		return newTemplateBindingsCommon(data.Common), nil
	default:
		return nil, fmt.Errorf("unexpected cursor ParsedData type %T for run", parsed)
	}
}
