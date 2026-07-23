package runtime

import "testing"

func TestDetectPayloadWith(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		raw    string
		getenv func(string) string
		want   bool
	}{
		{
			name: "cursor_version field",
			raw:  `{"cursor_version":"1.0","hook_event_name":"sessionStart"}`,
			want: true,
		},
		{
			name: "conversation_id field",
			raw:  `{"conversation_id":"c1","hook_event_name":"sessionStart"}`,
			want: true,
		},
		{
			name: "CURSOR_VERSION env",
			raw:  `{"hook_event_name":"sessionStart"}`,
			getenv: func(key string) string {
				if key == "CURSOR_VERSION" {
					return "1.2.3"
				}
				return ""
			},
			want: true,
		},
		{
			name: "no cursor markers",
			raw:  `{"session_id":"s","hook_event_name":"SessionStart"}`,
			getenv: func(string) string {
				return ""
			},
			want: false,
		},
		{
			name: "invalid json",
			raw:  `{`,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := detectPayloadWith([]byte(tt.raw), tt.getenv); got != tt.want {
				t.Fatalf("detectPayloadWith = %v, want %v", got, tt.want)
			}
		})
	}
}
