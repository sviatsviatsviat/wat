package event

// PreferDurationField returns the Hooks-documented duration when present on the
// wire (including explicit 0), otherwise durationMs.
//
// Callers that must honor explicit duration: 0 should set durationPresent from
// the DecodeEvent after-callback (for example via hookkit.RawObjectField), not
// via a custom Event.UnmarshalJSON.
func PreferDurationField(duration, durationMs int64, durationPresent bool) int64 {
	if durationPresent {
		return duration
	}
	return durationMs
}
