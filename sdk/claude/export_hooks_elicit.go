package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/elicit/elicitation"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/elicit/elicitationresult"
)

// Elicitation is the Elicitation hook event.
type Elicitation = elicitation.Event

// ElicitationOutput is the response for Elicitation events.
type ElicitationOutput = elicitation.Output

// ElicitationResults is the hook-scoped response builder for Elicitation.
type ElicitationResults = elicitation.Results

// ElicitationResult is the ElicitationResult hook event.
type ElicitationResult = elicitationresult.Event
