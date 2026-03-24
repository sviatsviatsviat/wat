package cursorcore

import (
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/core"
)

// DefaultHookResponseLine is the Cursor hook stdout line (JSON object and newline).
const DefaultHookResponseLine = "{}\n"

// EventHookHandler runs cmd with event template bindings.
type EventHookHandler struct {
	templateBindings core.TemplateBindings
}

// NewEventHookHandler returns a [core.HookHandler] that only wires [core.TemplateBindings] and runs cmd.Execute.
func NewEventHookHandler(templateBindings core.TemplateBindings) *EventHookHandler {
	return &EventHookHandler{templateBindings: templateBindings}
}

// NewEventHookHandlerFromExtractedFields builds an [EventHookHandler] from already-parsed common data and
// event-specific fields (constructor for callers that parse JSON themselves, e.g. decorators).
func NewEventHookHandlerFromExtractedFields[T any](
	commonData HookDataCommon,
	eventFields T,
	extractors map[string]EventFieldExtractor[T],
) *EventHookHandler {
	bindings := NewTemplateBindingsEvent(commonData, eventFields, extractors)
	return NewEventHookHandler(bindings)
}

// Handle runs cmd with [core.TemplateBindings] and fixed stdout payload.
func (handler *EventHookHandler) Handle(cmd core.Command) core.HookHandlerResult {
	ctx := &core.HookContext{
		TemplateBindings: handler.templateBindings,
	}
	code := cmd.Execute(ctx)
	return core.HookHandlerResult{Code: code, Output: DefaultHookResponseLine}
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
	return NewEventHookHandlerFromExtractedFields(hookData.HookDataCommon, hookData.Fields, extractors), nil
}

// HookHandlerBuilder builds a [core.HookHandler] from parsed hook fields and [core.WatExecutionContext].
type HookHandlerBuilder func(rawJSON []byte, hookData HookDataCommon, execCtx core.WatExecutionContext) (core.HookHandler, error)

// NewCursorEventHookHandlerBuilder returns a [HookHandlerBuilder] that parses event-specific JSON fields
// and wires [core.TemplateBindings] from extractors plus common placeholders.
func NewCursorEventHookHandlerBuilder[T any](extractors map[string]EventFieldExtractor[T]) HookHandlerBuilder {
	return func(rawJSON []byte, hookData HookDataCommon, _ core.WatExecutionContext) (core.HookHandler, error) {
		return newEventHookHandlerFromFields(rawJSON, hookData, extractors)
	}
}
