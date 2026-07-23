package prompt

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/prompt/beforesubmitprompt"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	beforesubmitprompt.Register(c)
}
