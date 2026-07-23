package tool

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/permissionrequest"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool/pretooluse"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	pretooluse.Register(c)
	posttooluse.Register(c)
	posttoolusefailure.Register(c)
	permissionrequest.Register(c)
}
