package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/prompt/beforesubmitprompt"
)

// BeforeSubmitPrompt is the beforeSubmitPrompt hook event.
// BeforeSubmitPrompt hooks.json matchers use the value UserPromptSubmit.
type BeforeSubmitPrompt = beforesubmitprompt.Event

// BeforeSubmitPromptOutput is the response for beforeSubmitPrompt events.
// BeforeSubmitPromptOutput from Block encodes continue:false with process
// exit 0 (not exit 2).
type BeforeSubmitPromptOutput = beforesubmitprompt.Output

// BeforeSubmitPromptResults is the hook-scoped response builder for BeforeSubmitPrompt.
// BeforeSubmitPromptResults.Block writes continue:false and user_message with
// exit 0.
type BeforeSubmitPromptResults = beforesubmitprompt.Results
