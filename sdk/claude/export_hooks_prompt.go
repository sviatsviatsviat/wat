package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt/userpromptexpansion"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt/userpromptsubmit"
)

// UserPromptSubmit is the UserPromptSubmit hook event.
type UserPromptSubmit = userpromptsubmit.Event

// UserPromptSubmitOutput is the response for UserPromptSubmit events.
type UserPromptSubmitOutput = userpromptsubmit.Output

// UserPromptSubmitResults is the hook-scoped response builder for UserPromptSubmit.
type UserPromptSubmitResults = userpromptsubmit.Results

// UserPromptExpansion is the UserPromptExpansion hook event.
type UserPromptExpansion = userpromptexpansion.Event

// UserPromptExpansionResults is the hook-scoped response builder for UserPromptExpansion.
type UserPromptExpansionResults = userpromptexpansion.Results
