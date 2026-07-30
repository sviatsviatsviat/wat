package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

// Dialect is the sdk/run router name for Claude Code hooks.
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
// Claude treats exit 1 as fail-open for most events; WorktreeCreate fails on any non-zero exit.
const HandlerErrorExit = event.HandlerErrorExit

// BlockExit is exit code 2 for TeammateIdle / TaskCreated / TaskCompleted Block builders.
// Claude feeds stderr to the model and ignores stdout JSON. Other Claude denies use
// exit 0 with JSON decision fields instead.
const BlockExit = event.BlockExit
