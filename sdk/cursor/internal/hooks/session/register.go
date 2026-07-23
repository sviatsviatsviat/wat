package session

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/sessionend"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/sessionstart"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session/workspaceopen"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	sessionstart.Register(c)
	sessionend.Register(c)
	workspaceopen.Register(c)
}
