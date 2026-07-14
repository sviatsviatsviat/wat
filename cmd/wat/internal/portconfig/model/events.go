package model

import "github.com/sviatsviatsviat/wat/sdk/agnostic"

// EventNameForEmit picks the native event name to use when emitting a handler entry.
func EventNameForEmit(e Entry, kindForEvent map[string]agnostic.Kind, eventForKind map[agnostic.Kind]string) string {
	if e.NativeEvent != "" {
		if k, ok := kindForEvent[e.NativeEvent]; ok && k == e.Kind {
			return e.NativeEvent
		}
	}
	return eventForKind[e.Kind]
}
