package permissionrequest

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/run"
)

// Output is the response for PermissionRequest events.
// Construct via Results builders and With* methods. A nil value is a no-op.
type Output interface {
	run.Output
	isOutput()
	// WithInterrupt stops the session when true.
	WithInterrupt(v bool) Output
}

type output struct {
	behavior         string
	message          string
	interrupt        bool
	suppressWarnExit bool
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.behavior == "" && o.message == "" && !o.interrupt
}

// WithInterrupt stops the session when true.
func (o output) WithInterrupt(v bool) Output {
	o.interrupt = v
	return o
}

// Encode renders this output as Copilot stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	if o.behavior == "" && o.message == "" && !o.interrupt {
		return nil, 0, nil
	}
	out := map[string]any{}
	if o.behavior != "" {
		out["behavior"] = o.behavior
	}
	if o.message != "" {
		out["message"] = o.message
	}
	if o.interrupt {
		out["interrupt"] = true
	}
	b, err := json.Marshal(out)
	if err != nil {
		return nil, 0, err
	}
	exitCode := 0
	if o.behavior == "deny" && !o.suppressWarnExit {
		exitCode = event.WarnExit
	}
	return b, exitCode, err
}

// Merge combines other into the receiver. other must be an output.
func (o output) Merge(other run.Output) (run.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	behavior, message := hookkit.MergeRankedString(
		o.behavior, o.message,
		b.behavior, b.message,
		hookkit.PermissionRankString,
	)
	oHardDeny := o.behavior == "deny" && !o.suppressWarnExit
	bHardDeny := b.behavior == "deny" && !b.suppressWarnExit
	suppressWarnExit := false
	if behavior == "deny" && !oHardDeny && !bHardDeny {
		suppressWarnExit = o.suppressWarnExit || b.suppressWarnExit
	}
	return output{
		behavior:         behavior,
		message:          message,
		interrupt:        o.interrupt || b.interrupt,
		suppressWarnExit: suppressWarnExit,
	}, nil, nil
}

// Stop reports whether remaining handlers should be skipped.
// Ask (deny with suppressWarnExit) does not stop.
func (o output) Stop() bool {
	return o.behavior == "deny" && !o.suppressWarnExit
}
