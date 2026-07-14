package hookkit

import "strings"

// ApplyRankedDecision merges a ranked decision field and optional paired detail
// (reason, message, etc.) from src into dst. When the new rank exceeds the old
// rank without a detail in src, any existing detail is cleared.
func ApplyRankedDecision(dst, src map[string]any, decisionKey, detailKey string, value any, rank func(any) int) {
	if rank(value) < rank(dst[decisionKey]) {
		return
	}
	oldRank := rank(dst[decisionKey])
	newRank := rank(value)
	dst[decisionKey] = value
	if detail, ok := src[detailKey]; ok && newRank > 0 {
		dst[detailKey] = detail
	} else if newRank > oldRank {
		delete(dst, detailKey)
	}
}

// ApplyOrphanDetail copies detailKey when dst has no decision yet.
func ApplyOrphanDetail(dst map[string]any, decisionKey, detailKey string, value any) {
	if value == nil {
		return
	}
	if _, has := dst[decisionKey]; !has && value != "" {
		dst[detailKey] = value
	}
}

// PermissionRank returns deny > ask > allow > unknown for permission strings.
func PermissionRank(v any) int {
	s, _ := v.(string)
	switch strings.ToLower(s) {
	case "deny":
		return 3
	case "ask":
		return 2
	case "allow":
		return 1
	default:
		return 0
	}
}

// BlockDecisionRank returns 1 for block decisions, else 0.
func BlockDecisionRank(v any) int {
	s, _ := v.(string)
	if strings.EqualFold(s, "block") {
		return 1
	}
	return 0
}

// JoinContext concatenates additional-context strings with a blank line separator.
func JoinContext(existing, add any) string {
	a, _ := existing.(string)
	b, _ := add.(string)
	switch {
	case a == "":
		return b
	case b == "":
		return a
	default:
		return a + "\n\n" + b
	}
}
