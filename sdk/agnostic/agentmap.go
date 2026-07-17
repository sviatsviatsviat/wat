package agnostic

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	sdkclaude "github.com/sviatsviatsviat/wat/sdk/claude"
	sdkcopilot "github.com/sviatsviatsviat/wat/sdk/copilot"
	sdkcursor "github.com/sviatsviatsviat/wat/sdk/cursor"
)

func claudeEvent(native sdkclaude.Event, raw []byte, kind Kind) *Event {
	env := sdkclaude.EnvelopeOf(native)
	return &Event{
		Agent:          Claude,
		Kind:           kind,
		Name:           native.EventName(),
		Session:        env.SessionID,
		Cwd:            env.Cwd,
		TranscriptPath: env.TranscriptPath,
		Raw:            hookkit.CloneRaw(raw),
	}
}

func copilotEvent(native sdkcopilot.Event, raw []byte, kind Kind) *Event {
	env := sdkcopilot.EnvelopeOf(native)
	name := native.EventName()
	if received := env.ReceivedName(); received != "" {
		name = received
	}
	return &Event{
		Agent:          Copilot,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            hookkit.CloneRaw(raw),
	}
}

func cursorEvent(native sdkcursor.Event, raw []byte, kind Kind) *Event {
	env := sdkcursor.EnvelopeOf(native)
	name := native.EventName()
	if received := env.ReceivedName(); received != "" {
		name = received
	}
	return &Event{
		Agent:          Cursor,
		Kind:           kind,
		Name:           name,
		Session:        env.Session(),
		Cwd:            env.Cwd,
		TranscriptPath: env.Transcript(),
		Raw:            hookkit.CloneRaw(raw),
	}
}

func cursorReceivedName(native sdkcursor.Event) string {
	if name := sdkcursor.EnvelopeOf(native).ReceivedName(); name != "" {
		return name
	}
	return native.EventName()
}
