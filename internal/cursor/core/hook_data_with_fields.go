package cursorcore

import "encoding/json"

// HookDataWithCommon composes shared Cursor hook metadata with event-specific fields.
type HookDataWithCommon[T any] struct {
	HookDataCommon
	Fields T
}

// NewHookDataWithCommon parses event-specific fields and composes them with commonData.
func NewHookDataWithCommon[T any](rawJSON []byte, commonData HookDataCommon) (HookDataWithCommon[T], error) {
	var fields T
	if err := json.Unmarshal(rawJSON, &fields); err != nil {
		return HookDataWithCommon[T]{}, err
	}
	return HookDataWithCommon[T]{
		HookDataCommon: commonData,
		Fields:         fields,
	}, nil
}
