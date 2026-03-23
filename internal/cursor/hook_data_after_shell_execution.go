package cursor

import "encoding/json"

// hookDataAfterShellExecution is the Cursor payload shape for afterShellExecution.
type hookDataAfterShellExecution struct {
	hookDataCommon
	hookDataAfterShellExecutionFields
}

type hookDataAfterShellExecutionFields struct {
	Command  string `json:"command"`
	Output   string `json:"output"`
	Duration float32 `json:"duration"`
	Sandbox  bool   `json:"sandbox"`
}

// newHookDataAfterShellExecution unmarshals only event-specific fields and composes them
// with already-parsed shared hookDataCommon from the factory path.
func newHookDataAfterShellExecution(rawJSON []byte, commonData hookDataCommon) (hookDataAfterShellExecution, error) {
	var fields hookDataAfterShellExecutionFields
	if err := json.Unmarshal(rawJSON, &fields); err != nil {
		return hookDataAfterShellExecution{}, err
	}
	return hookDataAfterShellExecution{
		hookDataCommon:                   commonData,
		hookDataAfterShellExecutionFields: fields,
	}, nil
}
