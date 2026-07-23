package stop

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop/stopfailure"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	stopevent.Register(c)
	stopfailure.Register(c)
}
