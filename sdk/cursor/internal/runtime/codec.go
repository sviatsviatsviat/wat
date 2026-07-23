package runtime

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Codec decodes Cursor hook stdin payloads.
var Codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)
