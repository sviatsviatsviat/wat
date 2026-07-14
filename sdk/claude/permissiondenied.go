package claude

import (
	"encoding/json"
)

// PermissionDenied is the PermissionDenied hook event.
type PermissionDenied struct {
	Envelope
	// ToolName is the tool name.
	ToolName string `json:"tool_name"`
	// ToolInput is the native tool input JSON.
	ToolInput json.RawMessage `json:"tool_input"`
}

// EventName returns the hook event name.
func (PermissionDenied) EventName() string { return EventPermissionDenied }

func init() {
	registerDecoder(EventPermissionDenied, decodeAs[PermissionDenied])
}

// PermissionDeniedOutput is the response for PermissionDenied events.
type PermissionDeniedOutput struct {
	Common
	// Retry requests a permission retry when true.
	Retry bool
}

func (o PermissionDeniedOutput) isZero() bool {
	return o.Common.isZero() && !o.Retry
}
