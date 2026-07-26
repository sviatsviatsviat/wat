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
func (d DurationFields) DurationMillis() int64 {
	return PreferDurationField(d.Duration, d.DurationMs, d.present)
}

// PreferDurationField returns the Hooks-documented duration when present on the
// wire (including explicit 0), otherwise durationMs.
//
// Callers that must honor explicit duration: 0 should set durationPresent from
// the DecodeEvent after-callback (for example via CaptureDurationPresent), not
// via a custom Event.UnmarshalJSON.
func PreferDurationField(duration, durationMs int64, durationPresent bool) int64 {
	if durationPresent {
		return duration
	}
	return durationMs
}
