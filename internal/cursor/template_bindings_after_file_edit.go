package cursor

import "github.com/sviatsviatsviat/wat/internal/cursor/core"

var afterFileEditPlaceholderExtractors = map[string]cursorcore.EventFieldExtractor[hookDataAfterFileEditFields]{
	"FILE_PATH": func(hookData hookDataAfterFileEditFields) string { return hookData.FilePath },
}
