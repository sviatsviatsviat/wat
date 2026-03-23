package cursorcore

import "github.com/sviatsviatsviat/wat/internal/core"

// EventFieldExtractor returns one template field from event-specific hook payload fields.
type EventFieldExtractor[T any] func(T) string

// templateBindingsEvent composes event-specific extractors with common placeholder bindings.
type templateBindingsEvent[T any] struct {
	commonBindings core.TemplateBindings
	eventFields    T
	extractors     map[string]EventFieldExtractor[T]
}

// NewTemplateBindingsEvent returns [core.TemplateBindings] for common plus event-specific placeholders.
func NewTemplateBindingsEvent[T any](
	commonData HookDataCommon,
	eventFields T,
	extractors map[string]EventFieldExtractor[T],
) core.TemplateBindings {
	return templateBindingsEvent[T]{
		commonBindings: newTemplateBindingsCommon(commonData),
		eventFields:    eventFields,
		extractors:     extractors,
	}
}

// TemplateValue returns the bound string for placeholderKey and whether it is a known key.
func (bindings templateBindingsEvent[T]) TemplateValue(placeholderKey string) (string, bool) {
	extractString, found := bindings.extractors[placeholderKey]
	if found && extractString != nil {
		return extractString(bindings.eventFields), true
	}
	return bindings.commonBindings.TemplateValue(placeholderKey)
}
