package model

import (
	"context"
)

// PreCompactEvent is the normalized view of a PreCompact hook invocation.
type PreCompactEvent struct {
	Envelope
	Compact *CompactInfo
}

// PreCompactHandler handles observe-only PreCompact events.
type PreCompactHandler func(ctx context.Context, event PreCompactEvent) error
