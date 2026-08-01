package stopevent

// Results is the hook-scoped response builder for this event.
type Results interface {
	// Context injects additionalContext that continues the Claude turn
	// (same stop_hook_active loop protection as FollowUp). Gate on
	// StopHookActive before returning Context.
	Context(text string) Output
	// FollowUp blocks completion with decision "block" and feeds reason back
	// to Claude. Gate on StopHookActive before returning FollowUp.
	FollowUp(reason string) Output
	isResults()
}

type results struct {
	eventName string
}

func (results) isResults() {}

// Context injects additionalContext that continues the Claude turn.
// Gate on StopHookActive before returning Context.
func (r results) Context(text string) Output {
	return output{eventName: r.eventName, additionalContext: text}
}

// FollowUp blocks completion with decision "block" and feeds reason back to Claude.
// Gate on StopHookActive before returning FollowUp.
func (r results) FollowUp(reason string) Output {
	return output{eventName: r.eventName, block: true, reason: reason}
}
