package doctor

import (
	"testing"
)

func TestParseGoModDirective(t *testing.T) {
	got, err := ParseGoModDirective([]byte("module wat-hooks\n\ngo 1.26\n"))
	if err != nil {
		t.Fatal(err)
	}
	if got != "1.26" {
		t.Fatalf("got %q want 1.26", got)
	}
}

func TestGoVersionAtLeast(t *testing.T) {
	tests := []struct {
		installed string
		required  string
		want      bool
	}{
		{"1.26.0", "1.26", true},
		{"1.26.1", "1.26", true},
		{"1.25.9", "1.26", false},
		{"1.26", "1.26.0", false}, // language version is less than the .0 release
		{"1.26rc1", "1.26", true},
		{"1.26.0", "1.26rc1", true},
		{"1.22-b9a08f159d", "1.22", true},
	}
	for _, tt := range tests {
		if got := GoVersionAtLeast(tt.installed, tt.required); got != tt.want {
			t.Fatalf("GoVersionAtLeast(%q, %q) = %v, want %v", tt.installed, tt.required, got, tt.want)
		}
	}
}
