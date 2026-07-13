package hookkit

import (
	"fmt"
	"reflect"
)

// NormalizeOutput dereferences a pointer output value when non-nil.
func NormalizeOutput(out any) any {
	if out == nil {
		return nil
	}
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer {
		return out
	}
	if v.IsNil() {
		return nil
	}
	return v.Elem().Interface()
}

// IsZeroOutput reports whether out is zero using reflection only.
func IsZeroOutput(out any) bool {
	if out == nil {
		return true
	}
	return reflect.ValueOf(out).IsZero()
}

// ValidateEncodePair checks that eventName is allowed for the given output type.
// canonicalize may remap eventName before comparison; pass nil to compare directly.
func ValidateEncodePair(label, eventName string, out any, allowed []string, canonicalize func(string) (string, bool)) error {
	if eventName == "" {
		return nil
	}
	canonical := eventName
	if canonicalize != nil {
		if mapped, ok := canonicalize(eventName); ok {
			canonical = mapped
		}
	}
	for _, name := range allowed {
		if canonical == name {
			return nil
		}
	}
	return fmt.Errorf("%s: encode: event %q incompatible with output type %T", label, eventName, out)
}
