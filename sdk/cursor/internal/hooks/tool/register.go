package tool

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/afterfileedit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftermcpexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftershellexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/aftertabfileedit"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforemcpexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforereadfile"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforeshellexecution"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/beforetabfileread"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/posttooluse"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/posttoolusefailure"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool/pretooluse"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	pretooluse.Register(c)
	posttooluse.Register(c)
	posttoolusefailure.Register(c)
	beforeshellexecution.Register(c)
	aftershellexecution.Register(c)
	beforemcpexecution.Register(c)
	aftermcpexecution.Register(c)
	beforereadfile.Register(c)
	afterfileedit.Register(c)
	beforetabfileread.Register(c)
	aftertabfileedit.Register(c)
}
