package runtime

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Codec decodes GitHub Copilot hook stdin payloads.
var Codec = hookkit.NewCodec(Dialect, ErrEmptyPayload, ErrDecodePayload, ErrEventNameRequired)
