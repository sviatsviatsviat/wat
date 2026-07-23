package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/compact/postcompact"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/compact/precompact"
)

// PreCompact is the PreCompact hook event.
type PreCompact = precompact.Event

// PreCompactResults is the hook-scoped response builder for PreCompact.
type PreCompactResults = precompact.Results

// PostCompact is the PostCompact hook event.
type PostCompact = postcompact.Event
