package errors

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/copilot/internal/hooks/errors/erroroccurred"
)

// Register registers this domain's hook event decoders on c.
func Register(c *hookkit.Codec) {
	erroroccurred.Register(c)
}
