package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

// MapSubagentStop maps a Cursor SubagentStop hook into a unified Event.
func MapSubagentStop(e sdkcursor.SubagentStop, raw []byte) *model.Event {
	ev := newEvent(e, raw, model.KindSubagentStop)
	tp := ""
	if e.AgentTranscriptPath != nil {
		tp = *e.AgentTranscriptPath
	}
	ev.Subagent = &model.Subagent{
		ID:             e.SubagentID,
		Type:           e.SubagentType,
		Task:           e.Task,
		Summary:        e.Summary,
		Status:         e.Status,
		TranscriptPath: tp,
		LoopCount:      e.LoopCount,
	}
	return ev
}
