package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/prompt/userpromptsubmitted"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/prompt/userprompttransformed"
)

// UserPromptSubmitted is the UserPromptSubmitted hook event.
type UserPromptSubmitted = userpromptsubmitted.Event

// UserPromptTransformed is the UserPromptTransformed hook event.
type UserPromptTransformed = userprompttransformed.Event

// UserPromptTransformedOutput is the response for UserPromptTransformed events.
type UserPromptTransformedOutput = userprompttransformed.Output

// UserPromptTransformedResults is the hook-scoped response builder for UserPromptTransformed.
type UserPromptTransformedResults = userprompttransformed.Results
