package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/prompt/beforesubmitprompt"
)

// BeforeSubmitPrompt is the beforeSubmitPrompt hook event.
type BeforeSubmitPrompt = beforesubmitprompt.Event

// BeforeSubmitPromptOutput is the response for beforeSubmitPrompt events.
type BeforeSubmitPromptOutput = beforesubmitprompt.Output

// BeforeSubmitPromptResults is the hook-scoped response builder for BeforeSubmitPrompt.
type BeforeSubmitPromptResults = beforesubmitprompt.Results
