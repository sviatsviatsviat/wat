package claude

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// outputEncoder is implemented by concrete hook outputs for Encode.
type outputEncoder interface {
	isZero() bool
	allowedEvents() []string
	encodeInto(top, hso map[string]any)
}

// Encode renders a typed output as Claude Code stdout JSON.
// eventName is written into hookSpecificOutput.hookEventName.
// A nil or zero output produces no stdout.
func Encode(eventName string, out any, opts ...Option) ([]byte, error) {
	cfg := defaultRuntimeConfig()
	applyOptions(&cfg, opts...)
	if eventName == "" {
		return nil, fmt.Errorf("claude: encode: empty event name")
	}
	out = hookkit.NormalizeOutput(out)
	if err := applyEnvSideEffect(eventName, out, cfg); err != nil {
		return nil, err
	}
	if out == nil {
		return nil, nil
	}
	enc, ok := out.(outputEncoder)
	if !ok {
		return nil, fmt.Errorf("claude: encode: unsupported output type %T", out)
	}
	if enc.isZero() {
		return nil, nil
	}
	if err := hookkit.ValidateEncodePair("claude", eventName, out, enc.allowedEvents(), nil); err != nil {
		return nil, err
	}

	top, hso := map[string]any{}, map[string]any{}
	enc.encodeInto(top, hso)
	if len(hso) > 0 {
		hso["hookEventName"] = eventName
		top["hookSpecificOutput"] = hso
	}
	if len(top) == 0 {
		return nil, nil
	}
	return json.Marshal(top)
}

func applyEnvSideEffect(eventName string, out any, cfg runtimeConfig) error {
	if eventName != EventSessionStart {
		return nil
	}
	w, ok := out.(interface{ writeSessionEnv(runtimeConfig) error })
	if !ok {
		return nil
	}
	return w.writeSessionEnv(cfg)
}

func isZeroOutput(out any) bool {
	if z, ok := out.(interface{ isZero() bool }); ok {
		return z.isZero()
	}
	return hookkit.IsZeroOutput(out)
}

// IsZeroOutput reports whether out is an empty hook response.
func IsZeroOutput(out any) bool { return isZeroOutput(out) }
