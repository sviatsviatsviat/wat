package agnostic

import "testing"

func TestCodecFor(t *testing.T) {
	tests := []struct {
		name    string
		dialect Dialect
		wantErr bool
	}{
		{name: "claude", dialect: Claude},
		{name: "copilot", dialect: Copilot},
		{name: "cursor", dialect: Cursor},
		{name: "unknown", dialect: Unknown, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			codec, err := CodecFor(tt.dialect)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if codec.Dialect() != tt.dialect {
				t.Fatalf("Dialect() = %v, want %v", codec.Dialect(), tt.dialect)
			}
		})
	}
}
