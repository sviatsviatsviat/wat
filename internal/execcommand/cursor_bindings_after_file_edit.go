package execcommand

import "github.com/sviatsviatsviat/wat/internal/cursor"

var afterFileEditPlaceholderExtractors = map[string]eventFieldExtractor[cursor.AfterFileEditFields]{
	"FILE_PATH": func(hookData cursor.AfterFileEditFields) string { return hookData.FilePath },
}
