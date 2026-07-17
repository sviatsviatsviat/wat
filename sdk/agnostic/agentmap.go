package agnostic

import (
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/adapter"
	"github.com/sviatsviatsviat/wat/sdk/agnostic/internal/model"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func claudeEvent(native sdkclaude.Event, raw []byte, kind model.Kind) *model.Event {
	env := sdkclaude.EnvelopeOf(native)
	return &model.Event{
		Agent:          model.Claude,
		Kind:           kind,
		Name:           native.EventName(),
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
		Raw:            adapter.CloneRaw(raw),
	}
}

func copilotEvent(native sdkcopilot.Event, raw []byte, kind model.Kind) *model.Event {
	env := sdkcopilot.EnvelopeOf(native)
	name := native.EventName()
	if received := env.ReceivedName(); received != "" {
		name = received
	}
	return &model.Event{
		Agent:          model.Copilot,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}
}

func cursorEvent(native sdkcursor.Event, raw []byte, kind model.Kind) *model.Event {
	env := sdkcursor.EnvelopeOf(native)
	name := native.EventName()
	if received := env.ReceivedName(); received != "" {
		name = received
	}
	return &model.Event{
		Agent:          model.Cursor,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            adapter.CloneRaw(raw),
	}
}

func cursorReceivedName(native sdkcursor.Event) string {
	if name := sdkcursor.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}
