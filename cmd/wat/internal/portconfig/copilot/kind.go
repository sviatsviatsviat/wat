package copilot

import (
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	agcopilot "github.com/sviatsviatsviat/wat/sdk/agnostic/copilot"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
)

var kindForEventMap = buildKindForEvent()

func kindForEvent(name string) (agnostic.Kind, bool) {
	canonical, ok := sdkcopilot.CanonicalEventName(name)
	if !ok {
		return agnostic.KindOther, false
	}
	kind, ok := kindForEventMap[canonical]
	return kind, ok
}

func buildKindForEvent() map[string]agnostic.Kind {
	out := model.InvertEventForKind(agcopilot.EventForKind)
	for alias, canonical := range sdkcopilot.EventAliasMap() {
		if kind, ok := out[canonical]; ok {
			out[alias] = kind
		}
	}
	return out
}
