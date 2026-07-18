package run

import "fmt"

// Decode parses raw for a registered dialect using that dialect's Decode op.
// eventHint is passed through for payloads that omit the wire event name.
func Decode(dialect, eventHint string, raw []byte) (any, error) {
	ops, ok := defaultRegistry.dialectOps(dialect)
	if !ok {
		return nil, fmt.Errorf("run: unknown dialect %q", dialect)
	}
	if ops.Decode == nil {
		return nil, fmt.Errorf("run: %s: missing Decode", dialect)
	}
	return ops.Decode(raw, eventHint)
}
