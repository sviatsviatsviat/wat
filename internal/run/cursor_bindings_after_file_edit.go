package run

import cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"

var afterFileEditPlaceholderExtractors = map[string]eventFieldExtractor[cursorcore.AfterFileEditFields]{
	"FILE_PATH": func(hookData cursorcore.AfterFileEditFields) string { return hookData.FilePath },
}
