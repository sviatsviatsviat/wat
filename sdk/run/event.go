package run

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Event is implemented by every decoded per-agent hook event.
type Event = hookkit.Event

// Output is any hook response. Concrete per-agent types implement IsZero and Encode.
type Output = hookkit.Output
