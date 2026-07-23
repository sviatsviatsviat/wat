package agent

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/agentstop"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent/subagentstop"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	subagentstart.Register(c)
	subagentstop.Register(c)
	agentstop.Register(c)
}
