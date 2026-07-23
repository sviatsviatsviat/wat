package agent

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/subagentstop"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/taskcompleted"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/taskcreated"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent/teammateidle"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	subagentstart.Register(c)
	subagentstop.Register(c)
	taskcreated.Register(c)
	taskcompleted.Register(c)
	teammateidle.Register(c)
}
