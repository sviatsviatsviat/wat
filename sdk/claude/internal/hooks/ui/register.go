package ui

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui/messagedisplay"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui/notification"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	notification.Register(c)
	messagedisplay.Register(c)
}
