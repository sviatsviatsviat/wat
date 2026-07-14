package claude

// PostToolBatch is the PostToolBatch hook event.
type PostToolBatch struct {
	Envelope
}

// EventName returns the hook event name.
func (PostToolBatch) EventName() string { return EventPostToolBatch }

func init() {
	registerDecoder(EventPostToolBatch, decodeAs[PostToolBatch])
}
