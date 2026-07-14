package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func mapSubagentStop(e sdkcursor.SubagentStop, ev *model.Event) {
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
}
