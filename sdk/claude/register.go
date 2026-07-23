package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/agent"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/compact"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/elicit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/prompt"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/session"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/stop"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/tool"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/ui"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/runtime"
)

func init() {
	c := runtime.Codec
	session.Register(c)
	prompt.Register(c)
	tool.Register(c)
	agent.Register(c)
	stop.Register(c)
	ui.Register(c)
	workspace.Register(c)
	compact.Register(c)
	elicit.Register(c)
}
