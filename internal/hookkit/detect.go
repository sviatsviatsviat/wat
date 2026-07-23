package hookkit

// DetectFunc reports whether raw matches a dialect.
type DetectFunc func(raw []byte) bool
