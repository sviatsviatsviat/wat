package agent

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/afteragentresponse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/afteragentthought"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/stopevent"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/subagentstart"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent/subagentstop"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	stopevent.Register(c)
	subagentstart.Register(c)
	subagentstop.Register(c)
	afteragentresponse.Register(c)
	afteragentthought.Register(c)
}
