package workspace

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/configchange"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/cwdchanged"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/filechanged"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/instructionsloaded"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/worktreecreate"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/hooks/workspace/worktreeremove"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	cwdchanged.Register(c)
	filechanged.Register(c)
	worktreecreate.Register(c)
	worktreeremove.Register(c)
	instructionsloaded.Register(c)
	configchange.Register(c)
}
