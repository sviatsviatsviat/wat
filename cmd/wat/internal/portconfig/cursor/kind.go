package cursor

import (
	"github.com/sviatsviatsviat/wat/cmd/wat/internal/portconfig/model"
	"github.com/sviatsviatsviat/wat/sdk/agnostic"
	agcursor "github.com/sviatsviatsviat/wat/sdk/agnostic/cursor"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

var kindForEventMap = buildKindForEvent()

func kindForEvent(name string) (agnostic.Kind, bool) {
	kind, ok := kindForEventMap[name]
	return kind, ok
}

func buildKindForEvent() map[string]agnostic.Kind {
	out := model.InvertEventForKind(agcursor.EventForKind)
	for event, kind := range map[string]agnostic.Kind{
		sdkcursor.EventBeforeShellExecution: agnostic.KindPreTool,
		sdkcursor.EventAfterShellExecution:  agnostic.KindPostTool,
		sdkcursor.EventBeforeMCPExecution:   agnostic.KindPreTool,
		sdkcursor.EventAfterMCPExecution:    agnostic.KindPostTool,
		sdkcursor.EventBeforeReadFile:       agnostic.KindPreTool,
		sdkcursor.EventAfterFileEdit:        agnostic.KindPostTool,
		sdkcursor.EventAfterAgentResponse:   agnostic.KindOther,
		sdkcursor.EventAfterAgentThought:    agnostic.KindOther,
		sdkcursor.EventBeforeTabFileRead:    agnostic.KindOther,
		sdkcursor.EventAfterTabFileEdit:     agnostic.KindOther,
		sdkcursor.EventWorkspaceOpen:        agnostic.KindOther,
	} {
		out[event] = kind
	}
	return out
}
