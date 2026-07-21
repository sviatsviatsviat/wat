package hookkit

// Output is a hook response. Concrete per-agent types implement IsZero and Encode.
type Output interface {
	IsZero() bool
	Encode() (stdout []byte, exit int, err error)
}
