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

// OnAfterTabFileEdit registers an observe-only afterTabFileEdit handler.
func OnAfterTabFileEdit(fn func(context.Context, Hook[AfterTabFileEdit]) error) *chain {
	return (&chain{}).AfterTabFileEdit(fn)
}

// AfterTabFileEdit registers another AfterTabFileEdit handler on the chain.
func (c *chain) AfterTabFileEdit(fn func(context.Context, Hook[AfterTabFileEdit]) error) *chain {
	registerObserveHandler(fn)
	return c
}
