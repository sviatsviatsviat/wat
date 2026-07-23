package compact

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/compact/precompact"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	precompact.Register(c)
}
