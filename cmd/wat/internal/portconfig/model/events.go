package model

import "github.com/sviatsviatsviat/wat/sdk/agnostic"

// InvertEventForKind builds an event-name to kind map from a kind to event map.
func InvertEventForKind(m map[agnostic.Kind]string) map[string]agnostic.Kind {
	out := make(map[string]agnostic.Kind, len(m))
	for kind, event := range m {
		out[event] = kind
	}
	return out
}

// EventNameForEmit picks the native event name to use when emitting a handler entry.
func EventNameForEmit(e Entry, kindForEvent map[string]agnostic.Kind, eventForKind map[agnostic.Kind]string) string {
	if e.NativeEvent != "" {
		if k, ok := kindForEvent[e.NativeEvent]; ok && k == e.Kind {
			return e.NativeEvent
		}
	}
	return eventForKind[e.Kind]
}
