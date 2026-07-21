package run

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Event is implemented by every decoded per-agent hook event.
type Event = hookkit.Event

// Output is a hook response. Concrete per-agent types implement IsZero, Encode,
// Merge, and Stop.
type Output = hookkit.Output
