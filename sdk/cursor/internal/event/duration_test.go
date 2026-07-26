package event

import (
	"encoding/json"
	"testing"
)

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
