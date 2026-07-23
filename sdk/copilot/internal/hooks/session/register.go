package session

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/session/sessionstart"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	sessionstart.Register(c)
	sessionend.Register(c)
}
