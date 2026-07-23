package run

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Event is implemented by every decoded per-agent hook event.
type Event = hookkit.Event

// Output is a hook response. Concrete per-agent types implement IsZero, Encode,
// Merge, and Stop.
type Output = hookkit.Output

// Hook is the handler context for a typed hook event.
type Hook[E Event] = hookkit.Hook[E]

// Invocation exposes serve-time settings to hook handlers.
type Invocation = hookkit.Invocation

// Config holds resolved serve-time settings passed to handlers via context.
type Config = hookkit.Config
