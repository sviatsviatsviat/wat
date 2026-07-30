package copilot_test

import (
	"context"
	"testing"

	"github.com/sviatsviatsviat/wat/sdk/copilot"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

func TestUseHooks_UserPromptTransformedInspect(t *testing.T) {
	hooks := copilot.UseHooks().UserPromptTransformed(func(context.Context, copilot.UserPromptTransformed, copilot.UserPromptTransformedResults) (copilot.UserPromptTransformedOutput, error) {
		return nil, nil
	})
	manifest := run.Inspect(hooks)
	got := manifest.EventsFor(copilot.Dialect)
	if len(got) != 1 || got[0] != copilot.EventUserPromptTransformed {
		t.Fatalf("EventsFor() = %v, want [%s]", got, copilot.EventUserPromptTransformed)
	}
	var handlers int
	for _, reg := range manifest.Registrations {
		if reg.Dialect == copilot.Dialect && reg.Event == copilot.EventUserPromptTransformed {
			handlers = reg.HandlerCount
		}
	}
	if handlers != 1 {
		t.Fatalf("HandlerCount = %d, want 1", handlers)
	}
}
