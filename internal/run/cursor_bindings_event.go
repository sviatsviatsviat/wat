package run

import (
	"errors"

	cursorcore "github.com/sviatsviatsviat/wat/internal/cursor/core"
	"github.com/sviatsviatsviat/wat/internal/template"
)

type eventFieldExtractor[T any] func(T) string

type templateBindingsEvent[T any] struct {
	commonBindings template.TemplateBindings
	eventFields    T
	extractors     map[string]eventFieldExtractor[T]
}

func newTemplateBindingsEvent[T any](
	commonData cursorcore.HookDataCommon,
	eventFields T,
	extractors map[string]eventFieldExtractor[T],
) template.TemplateBindings {
	return &templateBindingsEvent[T]{
		commonBindings: newTemplateBindingsCommon(commonData),
		eventFields:    eventFields,
		extractors:     extractors,
	}
}

func templateBindingsFromCursorEventPayload[T any](
	data *cursorcore.CursorHookRunData[T],
	extractors map[string]eventFieldExtractor[T],
	nilEventSpecificErr string,
) (template.TemplateBindings, error) {
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
