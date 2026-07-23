package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

// Dialect is the sdk/run registry name for Claude Code hooks.
const Dialect = runtime.Dialect

var (
	// ErrEmptyPayload indicates stdin was empty.
	ErrEmptyPayload = runtime.ErrEmptyPayload
	// ErrDecodePayload indicates JSON parsing of the hook payload failed.
	ErrDecodePayload = runtime.ErrDecodePayload
	// ErrEventNameRequired indicates the payload omitted hook_event_name.
	ErrEventNameRequired = runtime.ErrEventNameRequired
)

// SuccessExit is exit code 0 for successful encode and no-op responses.
const SuccessExit = event.SuccessExit

// HandlerErrorExit is exit code 1 for mux processing failures (read/decode/handler/encode/write errors).
const HandlerErrorExit = event.HandlerErrorExit
