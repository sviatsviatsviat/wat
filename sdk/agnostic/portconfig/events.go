package portconfig

import "github.com/sviatsviatsviat/wat/sdk/agnostic"

func eventNameForEmit(e Entry, kindForEvent map[string]agnostic.Kind, eventForKind map[agnostic.Kind]string) string {
	if e.NativeEvent != "" {
		if k, ok := kindForEvent[e.NativeEvent]; ok && k == e.Kind {
			return e.NativeEvent
		}
	}
	return eventForKind[e.Kind]
}
