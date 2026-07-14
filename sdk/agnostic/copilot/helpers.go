package copilot

import sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"

func nativeReceivedName(native sdkcopilot.Event) string {
	if name := sdkcopilot.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}
