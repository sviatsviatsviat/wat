package copilot

import (
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/agent"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/compact"
	hookerrors "github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/errors"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/prompt"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/session"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/tool"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/ui"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/runtime"
)

func init() {
	c := runtime.Codec
	session.Register(c)
	prompt.Register(c)
	tool.Register(c)
	agent.Register(c)
	compact.Register(c)
	ui.Register(c)
	hookerrors.Register(c)
}
