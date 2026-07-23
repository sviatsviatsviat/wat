package prompt

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt/userpromptexpansion"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt/userpromptsubmit"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	userpromptsubmit.Register(c)
	userpromptexpansion.Register(c)
}
