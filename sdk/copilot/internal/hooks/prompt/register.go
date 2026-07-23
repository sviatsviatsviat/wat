package prompt

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/prompt/userpromptsubmitted"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	userpromptsubmitted.Register(c)
}
