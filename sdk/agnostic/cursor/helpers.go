package cursor

import sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"

func receivedName(native sdkcursor.Event) string {
	if name := sdkcursor.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}
