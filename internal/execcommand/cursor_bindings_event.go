package execcommand

import "github.com/sviatsviatsviat/wat/internal/cursor"

type eventFieldExtractor[T any] func(T) string

type templateBindingsEvent[T any] struct {
	commonBindings templateBindings
	eventFields    T
	extractors     map[string]eventFieldExtractor[T]
}

func newTemplateBindingsEvent[T any](
	commonData cursor.HookDataCommon,
	eventFields T,
	extractors map[string]eventFieldExtractor[T],
) templateBindings {
	return &templateBindingsEvent[T]{
		commonBindings: newTemplateBindingsCommon(commonData),
		eventFields:    eventFields,
		extractors:     extractors,
	}
}

// templateBindingsFromCursorEventPayload builds bindings from common and event-specific fields.
// event must be non-nil: Cursor event hook adapters always populate it from parsed stdin.
func templateBindingsFromCursorEventPayload[T any](
	common cursor.HookDataCommon,
	event *T,
	extractors map[string]eventFieldExtractor[T],
) templateBindings {
	return newTemplateBindingsEvent(common, *event, extractors)
}

func (bindings *templateBindingsEvent[T]) TemplateValue(placeholderKey string) (string, bool) {
	extractString, found := bindings.extractors[placeholderKey]
	if found && extractString != nil {
		return extractString(bindings.eventFields), true
	}
	return bindings.commonBindings.TemplateValue(placeholderKey)
}
