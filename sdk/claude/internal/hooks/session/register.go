package session

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session/setup"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	sessionstart.Register(c)
	sessionend.Register(c)
	setup.Register(c)
}
