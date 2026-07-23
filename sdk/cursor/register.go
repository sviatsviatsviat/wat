package cursor

import (
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/agent"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/compact"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/prompt"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/session"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/hooks/tool"
	"github.com/sviatsviatsviat/wat/sdk/cursor/internal/runtime"
)

func init() {
	c := runtime.Codec
	session.Register(c)
	prompt.Register(c)
	tool.Register(c)
	agent.Register(c)
	compact.Register(c)
	runtime.EnsureDialect(runtime.DefaultReg)
}
