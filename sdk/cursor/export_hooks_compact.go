package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/compact/precompact"
)

// PreCompact is the preCompact hook event, including Cursor compaction metrics
// (context usage, token counts, message counts, and first-compaction flag).
type PreCompact = precompact.Event

// PreCompactOutput is the response for preCompact events.
// Cursor documents this hook as observational; the only supported stdout field
// is an optional user_message (see PreCompactResults.UserMessage).
type PreCompactOutput = precompact.Output

// PreCompactResults is the hook-scoped response builder for PreCompact.
type PreCompactResults = precompact.Results
