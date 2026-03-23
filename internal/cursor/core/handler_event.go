package cursorcore

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// defaultHookResponseLine is the default hook stdout line (JSON object and newline).
const defaultHookResponseLine = "{}\n"

// EventHookHandler runs cmd with template bindings derived from event-specific hook data.
type EventHookHandler[T any] struct {
	templateBindings core.TemplateBindings
}

// Handle runs cmd with [core.TemplateBindings] and fixed stdout payload.
func (handler EventHookHandler[T]) Handle(cmd core.Command) core.HookHandlerResult {
	ctx := &core.HookContext{
		TemplateBindings: handler.templateBindings,
	}
	code := cmd.Execute(ctx)
	return core.HookHandlerResult{Code: code, Output: defaultHookResponseLine}
}

func newEventHookHandlerFromFields[T any](
	rawJSON []byte,
	commonData HookDataCommon,
	extractors map[string]EventFieldExtractor[T],
) (core.HookHandler, error) {
	hookData, err := NewHookDataWithCommon[T](rawJSON, commonData)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor %s payload: %w", commonData.HookEventName, err)
	}
	templateBindings := NewTemplateBindingsEvent(hookData.HookDataCommon, hookData.Fields, extractors)
	return EventHookHandler[HookDataWithCommon[T]]{
		templateBindings: templateBindings,
	}, nil
}

// HookHandlerBuilder builds a [core.HookHandler] from parsed hook fields.
type HookHandlerBuilder func(rawJSON []byte, hookData HookDataCommon) (core.HookHandler, error)

// NewCursorEventHookHandlerBuilder returns a [HookHandlerBuilder] that parses event-specific JSON fields
// and wires [core.TemplateBindings] from extractors plus common placeholders.
func NewCursorEventHookHandlerBuilder[T any](extractors map[string]EventFieldExtractor[T]) HookHandlerBuilder {
	return func(rawJSON []byte, hookData HookDataCommon) (core.HookHandler, error) {
		return newEventHookHandlerFromFields(rawJSON, hookData, extractors)
	}
}
