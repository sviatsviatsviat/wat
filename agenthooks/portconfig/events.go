package portconfig

import "github.com/sviatsviatsviat/wat/agenthooks"

func eventNameForEmit(e Entry, kindForEvent map[string]agenthooks.Kind, eventForKind map[agenthooks.Kind]string) string {
	if e.NativeEvent != "" {
		if k, ok := kindForEvent[e.NativeEvent]; ok && k == e.Kind {
			return e.NativeEvent
		}
	}
	return eventForKind[e.Kind]
}
