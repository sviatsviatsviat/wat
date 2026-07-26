package event

import (
	"encoding/json"
	"testing"
)

func TestPreferDurationField(t *testing.T) {
	cases := []struct {
		name            string
		duration        int64
		durationMs      int64
		durationPresent bool
		want            int64
	}{
		{name: "documented present", duration: 100, durationMs: 999, durationPresent: true, want: 100},
		{name: "explicit zero beats ms", duration: 0, durationMs: 999, durationPresent: true, want: 0},
		{name: "absent falls back to ms", duration: 0, durationMs: 50, durationPresent: false, want: 50},
		{name: "both absent", duration: 0, durationMs: 0, durationPresent: false, want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := PreferDurationField(tc.duration, tc.durationMs, tc.durationPresent); got != tc.want {
				t.Fatalf("PreferDurationField() = %d, want %d", got, tc.want)
			}
		})
	}
}

func TestDurationFields_DurationMillis(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want int64
	}{
		{name: "documented present", raw: `{"duration":100,"duration_ms":999}`, want: 100},
		{name: "explicit zero beats ms", raw: `{"duration":0,"duration_ms":999}`, want: 0},
		{name: "absent falls back to ms", raw: `{"duration_ms":50}`, want: 50},
		{name: "both absent", raw: `{}`, want: 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var d DurationFields
			if err := json.Unmarshal([]byte(tc.raw), &d); err != nil {
				t.Fatal(err)
			}
			d.CaptureDurationPresent([]byte(tc.raw))
			if got := d.DurationMillis(); got != tc.want {
				t.Fatalf("DurationMillis() = %d, want %d", got, tc.want)
			}
		})
	}
}
