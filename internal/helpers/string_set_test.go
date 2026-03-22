package helpers

import (
	"slices"
	"testing"
)

func TestStringSetSortedKeys(t *testing.T) {
	t.Parallel()

	t.Run("nil input", func(t *testing.T) {
		t.Parallel()
		if sortedKeys := StringSetSortedKeys(nil); sortedKeys != nil {
			t.Fatalf("nil map: want nil slice, got %#v", sortedKeys)
		}
	})

	t.Run("empty input", func(t *testing.T) {
		t.Parallel()
		if sortedKeys := StringSetSortedKeys(map[string]struct{}{}); sortedKeys != nil {
			t.Fatalf("empty set: want nil slice, got %#v", sortedKeys)
		}
	})

	t.Run("sorted keys", func(t *testing.T) {
		t.Parallel()
		sortedKeys := StringSetSortedKeys(map[string]struct{}{
			"z": {}, "a": {}, "m": {},
		})
		wantSorted := []string{"a", "m", "z"}
		if !slices.Equal(sortedKeys, wantSorted) {
			t.Fatalf("got %v, want %v", sortedKeys, wantSorted)
		}
	})
}
