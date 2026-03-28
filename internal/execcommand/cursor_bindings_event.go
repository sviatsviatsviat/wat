package execcommand

import (
	"errors"

	"github.com/sviatsviatsviat/wat/internal/cursor"
)

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

func templateBindingsFromCursorEventPayload[T any](
	data *cursor.CursorHookRunData[T],
	extractors map[string]eventFieldExtractor[T],
	nilEventSpecificErr string,
) (templateBindings, error) {
	if data.EventSpecific == nil {
		return nil, errors.New(nilEventSpecificErr)
	}
	return newTemplateBindingsEvent(data.Common, *data.EventSpecific, extractors), nil
}

func (bindings *templateBindingsEvent[T]) TemplateValue(placeholderKey string) (string, bool) {
	extractString, found := bindings.extractors[placeholderKey]
	if found && extractString != nil {
		return extractString(bindings.eventFields), true
	}
	return bindings.commonBindings.TemplateValue(placeholderKey)
}
