package event

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// DurationFields holds Cursor duration wire fields shared by post-tool events.
//
// Embed on events that carry `duration` / `duration_ms`, then call
// CaptureDurationPresent from the DecodeEvent after-callback so
// DurationMillis honors an explicit duration: 0.
type DurationFields struct {
	// Duration is the Hooks-documented duration field in milliseconds.
	Duration int64 `json:"duration"`
	// DurationMs is an alternate duration field in milliseconds.
	DurationMs int64 `json:"duration_ms"`
	present    bool
}

// CaptureDurationPresent records whether the documented duration key was
// present on raw JSON. Call from the DecodeEvent after-callback.
func (d *DurationFields) CaptureDurationPresent(raw []byte) {
	d.present = hookkit.RawObjectField(raw, "duration") != nil
}

// DurationMillis returns the Hooks-documented duration when present on the
// wire (including explicit 0), otherwise DurationMs.
// Prefer this helper over reading Duration or DurationMs directly.
//
// Call CaptureDurationPresent from the DecodeEvent after-callback so an
// explicit duration: 0 is distinguished from an absent duration key.
func (d DurationFields) DurationMillis() int64 {
	if d.present {
		return d.Duration
	}
	return d.DurationMs
}
