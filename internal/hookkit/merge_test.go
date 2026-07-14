package hookkit

import "testing"

func TestApplyRankedDecision(t *testing.T) {
	rank := PermissionRank

	tests := []struct {
		name string
		dst  map[string]any
		src  map[string]any
		want map[string]any
	}{
		{
			name: "stricter decision replaces with reason",
			dst:  map[string]any{"permissionDecision": "allow", "permissionDecisionReason": "ok"},
			src: map[string]any{
				"permissionDecision":       "deny",
				"permissionDecisionReason": "blocked",
			},
			want: map[string]any{"permissionDecision": "deny", "permissionDecisionReason": "blocked"},
		},
		{
			name: "stricter decision without reason clears stale reason",
			dst:  map[string]any{"permissionDecision": "allow", "permissionDecisionReason": "ok"},
			src:  map[string]any{"permissionDecision": "deny"},
			want: map[string]any{"permissionDecision": "deny"},
		},
		{
			name: "weaker decision ignored",
			dst:  map[string]any{"permissionDecision": "deny", "permissionDecisionReason": "blocked"},
			src:  map[string]any{"permissionDecision": "allow", "permissionDecisionReason": "ok"},
			want: map[string]any{"permissionDecision": "deny", "permissionDecisionReason": "blocked"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := mapsClone(tt.dst)
			ApplyRankedDecision(dst, tt.src, "permissionDecision", "permissionDecisionReason", tt.src["permissionDecision"], rank)
			if !mapsEqual(dst, tt.want) {
				t.Fatalf("dst = %v, want %v", dst, tt.want)
			}
		})
	}
}

func TestApplyOrphanDetail(t *testing.T) {
	dst := map[string]any{}
	ApplyOrphanDetail(dst, "permissionDecision", "permissionDecisionReason", "early reason")
	if dst["permissionDecisionReason"] != "early reason" {
		t.Fatalf("orphan detail = %v", dst)
	}

	dst = map[string]any{"permissionDecision": "deny"}
	ApplyOrphanDetail(dst, "permissionDecision", "permissionDecisionReason", "ignored")
	if _, ok := dst["permissionDecisionReason"]; ok {
		t.Fatalf("detail should not apply when decision exists: %v", dst)
	}
}

func TestJoinContext(t *testing.T) {
	if got := JoinContext("first", "second"); got != "first\n\nsecond" {
		t.Fatalf("JoinContext() = %q", got)
	}
}

func mapsClone(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func mapsEqual(a, b map[string]any) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
