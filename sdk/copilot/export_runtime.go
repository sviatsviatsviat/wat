package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

// Dialect is the sdk/run router name for GitHub Copilot hooks.
const Dialect = runtime.Dialect

var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = runtime.ErrEmptyPayload
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = runtime.ErrDecodePayload
	// ErrEventNameRequired indicates the payload omitted hook_event_name.
	ErrEventNameRequired = runtime.ErrEventNameRequired
)

// HandlerErrorExit is exit code 1 for handler errors. Copilot command hooks
// fail-closed on non-zero exits other than 2 (including PreToolUse).
const HandlerErrorExit = event.HandlerErrorExit

// WarnExit is exit code 2. Copilot treats it as a warning by default; for
// PermissionRequest it means deny, and for PostToolUseFailure it carries
// additional_context in stdout.
const WarnExit = event.WarnExit
