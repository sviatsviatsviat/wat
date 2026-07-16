package cursor

import (
	"context"
)

// AfterTabFileEdit is the afterTabFileEdit hook event.
type AfterTabFileEdit struct {
	Envelope
	// FilePath is the edited file path.
	FilePath string `json:"file_path"`
	// Edits are the applied edits.
	Edits []Edit `json:"edits"`
}

// EventName returns the canonical hook event name.
func (AfterTabFileEdit) EventName() string { return EventAfterTabFileEdit }

func init() {
	registerDecoder(EventAfterTabFileEdit, decodeAs[AfterTabFileEdit])
}

// AfterTabFileEdit registers an observe-only afterTabFileEdit handler.
func (c *Chain) AfterTabFileEdit(fn func(context.Context, Hook[AfterTabFileEdit]) error) *Chain {
	registerObserveHandler(fn)
	return &Chain{}
}
