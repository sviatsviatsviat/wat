package agnostic

import (
	"fmt"
	"strings"
)

// mappedDecodeError rewrites an SDK decode error for agnostic while preserving Unwrap.
type mappedDecodeError struct {
	msg string
	err error
}

func (e *mappedDecodeError) Error() string { return e.msg }

func (e *mappedDecodeError) Unwrap() error { return e.err }

func mapDecodeErrorMessage(err error, agent, sdk string) error {
	prefix := sdk + ": "
	msg := err.Error()
	if strings.HasPrefix(msg, prefix) {
		return &mappedDecodeError{
			msg: fmt.Sprintf("%s: %s", agent, msg[len(prefix):]),
			err: err,
		}
	}
	return fmt.Errorf("%s: %w", agent, err)
}

func remapDecodeError(err error, msg string) error {
	return &mappedDecodeError{msg: msg, err: err}
}
