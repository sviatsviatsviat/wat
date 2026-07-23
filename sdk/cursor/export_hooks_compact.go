package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/compact/precompact"
)

// PreCompact is the preCompact hook event.
type PreCompact = precompact.Event

// PreCompactOutput is the response for preCompact events.
type PreCompactOutput = precompact.Output

// PreCompactResults is the hook-scoped response builder for PreCompact.
type PreCompactResults = precompact.Results
