package workspaceopen

import (
	"encoding/json"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Output is the response for workspaceOpen events.
// Construct via Results builders and With* methods. A nil value is a no-op.
type Output interface {
	hookkit.Output
	isOutput()
	// WithPluginPaths sets absolute plugin directories to load for the workspace.
	WithPluginPaths(paths []string) Output
}

type output struct {
	pluginPaths []string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return len(o.pluginPaths) == 0
}

// WithPluginPaths sets absolute plugin directories to load for the workspace.
func (o output) WithPluginPaths(paths []string) Output {
	o.pluginPaths = append([]string(nil), paths...)
	return o
}

// Encode renders this output as Cursor stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	if len(o.pluginPaths) == 0 {
		return nil, 0, nil
	}
	out := map[string]any{"pluginPaths": o.pluginPaths}
	b, err := json.Marshal(out)
	return b, 0, err
}

// Merge combines other into this workspaceOpen output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	var warnings []string
	pluginPaths, w := hookkit.TakeLastSlice("pluginPaths", o.pluginPaths, b.pluginPaths)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{pluginPaths: pluginPaths}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return false
}
