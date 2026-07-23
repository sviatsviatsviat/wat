package claude

import (
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/config"
)

// Settings is the Claude Code settings.json hooks section shape.
type Settings = config.Settings

// MatcherGroup is a Claude Code hook matcher group.
type MatcherGroup = config.MatcherGroup

// Handler is a Claude Code hook handler definition.
type Handler = config.Handler
