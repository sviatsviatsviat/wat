package elicit

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/elicit/elicitation"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/elicit/elicitationresult"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	elicitation.Register(c)
	elicitationresult.Register(c)
}
