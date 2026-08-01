package hookkit

import (
	"fmt"
	"maps"
	"strings"
)

// OverwriteWarning returns a last-wins overwrite warning for field.
func OverwriteWarning(field string) string {
	return field + ": overwritten by later handler"
}

// TakeLastPtr keeps src when non-nil, else dst. When both are non-nil, returns
// src and an overwrite warning for field.
func TakeLastPtr[T any](field string, dst, src *T) (val *T, warning string) {
	if src == nil {
		return dst, ""
	}
	if dst != nil {
		return src, OverwriteWarning(field)
	}
	return src, ""
}

// TakeLastString keeps src when non-empty, else dst. When both are non-empty,
// returns src and an overwrite warning for field.
func TakeLastString(field, dst, src string) (val, warning string) {
	if src == "" {
		return dst, ""
	}
	if dst != "" {
		return src, OverwriteWarning(field)
	}
	return src, ""
}

// TakeLastMap keeps a clone of src when non-nil, else dst. When both are
// non-nil, returns a clone of src and an overwrite warning for field.
func TakeLastMap[K comparable, V any](field string, dst, src map[K]V) (val map[K]V, warning string) {
	if src == nil {
		return dst, ""
	}
	cloned := maps.Clone(src)
	if dst != nil {
		return cloned, OverwriteWarning(field)
	}
	return cloned, ""
}

// TakeLastSlice keeps a clone of src when non-nil, else dst. When both are
// non-nil, returns a clone of src and an overwrite warning for field.
func TakeLastSlice[T any](field string, dst, src []T) (val []T, warning string) {
	if src == nil {
		return dst, ""
	}
	cloned := append([]T(nil), src...)
	if dst != nil {
		return cloned, OverwriteWarning(field)
	}
	return cloned, ""
}

// CloneJSONValue returns a shallow clone of map/slice JSON payloads.
// Other values are returned as-is. Nil stays nil.
func CloneJSONValue(v any) any {
	if v == nil {
		return nil
	}
	switch x := v.(type) {
	case map[string]any:
		return maps.Clone(x)
	case []any:
		return append([]any(nil), x...)
	default:
		return v
	}
}

// TakeLastAny keeps a clone of src when non-nil, else dst. When both are
// non-nil, returns a clone of src and an overwrite warning for field.
// Map and slice payloads are cloned; other values are kept as-is.
func TakeLastAny(field string, dst, src any) (val any, warning string) {
	if src == nil {
		return dst, ""
	}
	cloned := CloneJSONValue(src)
	if dst != nil {
		return cloned, OverwriteWarning(field)
	}
	return cloned, ""
}

// MergeRankedString merges a ranked string decision and its paired detail.
// rank should return higher values for stricter decisions. When the incoming
// rank is lower, dst is unchanged. When the incoming rank is higher without a
// detail, the previous detail is cleared.
func MergeRankedString(dstDecision, dstDetail, srcDecision, srcDetail string, rank func(string) int) (decision, detail string) {
	if rank(srcDecision) < rank(dstDecision) {
		return dstDecision, dstDetail
	}
	oldRank := rank(dstDecision)
	newRank := rank(srcDecision)
	decision = srcDecision
	if srcDetail != "" {
		detail = srcDetail
	} else if newRank > oldRank {
		detail = ""
	} else {
		detail = dstDetail
	}
	return decision, detail
}

// PermissionRankString returns deny > defer > ask > allow > unknown for
// permission labels (Claude PreToolUse host mux order).
func PermissionRankString(s string) int {
	switch strings.ToLower(s) {
	case "deny":
		return 4
	case "defer":
		return 3
	case "ask":
		return 2
	case "allow":
		return 1
	default:
		return 0
	}
}

// BlockDecisionRankString returns 1 for block decisions, else 0.
func BlockDecisionRankString(s string) int {
	if strings.EqualFold(s, "block") {
		return 1
	}
	return 0
}

// ElicitationActionRankString returns decline/cancel > accept > unknown.
func ElicitationActionRankString(s string) int {
	switch strings.ToLower(s) {
	case "decline", "cancel":
		return 2
	case "accept":
		return 1
	default:
		return 0
	}
}

// JoinContextStrings concatenates additional-context strings with a blank line separator.
func JoinContextStrings(existing, add string) string {
	switch {
	case existing == "":
		return add
	case add == "":
		return existing
	default:
		return existing + "\n\n" + add
	}
}

// ErrMergeType returns an error when Merge receives a mismatched concrete type.
func ErrMergeType(want, got any) error {
	return fmt.Errorf("merge type mismatch: want %T, got %T", want, got)
}
