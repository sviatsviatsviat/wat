package hookkit

import (
	"encoding/json"
	"testing"
)

func TestTimestamp_UnmarshalJSON(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		raw      string
		want     int64
		wantErr  bool
		wantZero bool
	}{
		{name: "iso8601", raw: `"2026-07-12T10:00:00Z"`, want: 1783850400000},
		{name: "empty string", raw: `""`, wantZero: true},
		{name: "null", raw: `null`, wantZero: true},
		{name: "whitespace null", raw: `  null  `, wantZero: true},
		{name: "ms epoch rejected", raw: `1760000000000`, wantErr: true},
		{name: "invalid RFC3339", raw: `"not-a-time"`, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var ts Timestamp
			err := json.Unmarshal([]byte(tt.raw), &ts)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if tt.wantZero {
				if !ts.IsZero() {
					t.Fatalf("Timestamp=%v, want zero", ts)
				}
				return
			}
			if ts.UnixMilli() != tt.want {
				t.Fatalf("UnixMilli()=%d, want %d", ts.UnixMilli(), tt.want)
			}
		})
	}
}
