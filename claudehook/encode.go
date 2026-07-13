package claudehook

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Encode renders a typed output struct as Claude Code stdout JSON.
// eventName is written into hookSpecificOutput.hookEventName.
func Encode(eventName string, out any, opts ...Option) ([]byte, error) {
	if eventName == "" {
		return nil, fmt.Errorf("claudehook: encode: empty event name")
	}
	out = normalizeOutput(out)
	if out == nil || isZeroOutput(out) {
		return nil, nil
	}
	cfg := defaultRuntimeConfig()
	applyOptions(&cfg, opts...)

	top, hso, err := buildOutput(eventName, out)
	if err != nil {
		return nil, err
	}
	if err := applyEnvSideEffect(eventName, out, cfg); err != nil {
		return nil, err
	}
	if len(hso) > 0 {
		hso["hookEventName"] = eventName
		top["hookSpecificOutput"] = hso
	}
	if len(top) == 0 {
		return nil, nil
	}
	return json.Marshal(top)
}

func normalizeOutput(out any) any {
	if out == nil {
		return nil
	}
	v := reflect.ValueOf(out)
	if v.Kind() != reflect.Pointer {
		return out
	}
	if v.IsNil() {
		return nil
	}
	return v.Elem().Interface()
}

func isZeroOutput(out any) bool {
	if z, ok := out.(interface{ isZero() bool }); ok {
		return z.isZero()
	}
	return reflect.ValueOf(out).IsZero()
}

func buildOutput(eventName string, out any) (map[string]any, map[string]any, error) {
	top := map[string]any{}
	hso := map[string]any{}

	switch o := out.(type) {
	case PreToolUseOutput:
		applyCommon(top, o.Common)
		if o.Decision != "" {
			hso["permissionDecision"] = string(o.Decision)
			if o.Reason != "" {
				hso["permissionDecisionReason"] = o.Reason
			}
		} else if o.UpdatedInput != nil {
			hso["permissionDecision"] = "allow"
		}
		if o.UpdatedInput != nil {
			hso["updatedInput"] = o.UpdatedInput
		}
		if o.AdditionalContext != "" {
			hso["additionalContext"] = o.AdditionalContext
		}
	case PermissionRequestOutput:
		applyCommon(top, o.Common)
		if o.Behavior != "" {
			dec := map[string]any{"behavior": o.Behavior}
			if o.UpdatedInput != nil {
				dec["updatedInput"] = o.UpdatedInput
			}
			if o.Message != "" {
				dec["message"] = o.Message
			}
			if o.Interrupt {
				dec["interrupt"] = true
			}
			hso["decision"] = dec
		}
		if o.AdditionalContext != "" {
			hso["additionalContext"] = o.AdditionalContext
		}
	case PostToolUseOutput:
		applyCommon(top, o.Common)
		if o.Block {
			top["decision"] = "block"
			if o.Reason != "" {
				top["reason"] = o.Reason
			}
		}
		if o.UpdatedToolOutput != nil {
			hso["updatedToolOutput"] = o.UpdatedToolOutput
		}
		if o.AdditionalContext != "" {
			hso["additionalContext"] = o.AdditionalContext
		}
	case UserPromptSubmitOutput:
		applyCommon(top, o.Common)
		if o.Block {
			top["decision"] = "block"
			if o.Reason != "" {
				top["reason"] = o.Reason
			}
		}
		if o.AdditionalContext != "" {
			hso["additionalContext"] = o.AdditionalContext
		}
		if o.SessionTitle != "" {
			hso["sessionTitle"] = o.SessionTitle
		}
	case StopOutput:
		applyCommon(top, o.Common)
		if o.Block {
			top["decision"] = "block"
			if o.Reason != "" {
				top["reason"] = o.Reason
			}
		}
		if o.AdditionalContext != "" {
			hso["additionalContext"] = o.AdditionalContext
		}
	case SessionStartOutput:
		applyCommon(top, o.Common)
		if o.AdditionalContext != "" {
			hso["additionalContext"] = o.AdditionalContext
		}
		if o.SessionTitle != "" {
			hso["sessionTitle"] = o.SessionTitle
		}
		if o.InitialUserMessage != "" {
			hso["initialUserMessage"] = o.InitialUserMessage
		}
		if len(o.WatchPaths) > 0 {
			hso["watchPaths"] = o.WatchPaths
		}
		if o.ReloadSkills {
			hso["reloadSkills"] = true
		}
	case MessageDisplayOutput:
		applyCommon(top, o.Common)
		if o.DisplayContent != nil {
			hso["displayContent"] = *o.DisplayContent
		}
	case PermissionDeniedOutput:
		applyCommon(top, o.Common)
		if o.Retry {
			hso["retry"] = true
		}
	case ElicitationOutput:
		applyCommon(top, o.Common)
		if o.Action != "" {
			hso["action"] = o.Action
		}
		if o.Content != nil {
			hso["content"] = o.Content
		}
	case WorktreeCreateOutput:
		applyCommon(top, o.Common)
		if o.WorktreePath != "" {
			hso["worktreePath"] = o.WorktreePath
		}
	case CommonOutput:
		applyCommon(top, o.Common)
		if o.AdditionalContext != "" {
			hso["additionalContext"] = o.AdditionalContext
		}
	default:
		return nil, nil, fmt.Errorf("claudehook: encode: unsupported output type %T", out)
	}
	_ = eventName
	return top, hso, nil
}

func applyCommon(top map[string]any, c Common) {
	if c.Continue != nil && !*c.Continue {
		top["continue"] = false
		if c.StopReason != "" {
			top["stopReason"] = c.StopReason
		}
	}
	if c.SystemMessage != "" {
		top["systemMessage"] = c.SystemMessage
	}
	if c.SuppressOutput {
		top["suppressOutput"] = true
	}
	if c.TerminalSequence != "" {
		top["terminalSequence"] = c.TerminalSequence
	}
}

func applyEnvSideEffect(eventName string, out any, cfg runtimeConfig) error {
	if eventName != EventSessionStart {
		return nil
	}
	o, ok := out.(SessionStartOutput)
	if !ok || len(o.Env) == 0 {
		return nil
	}
	return WriteEnvFile(o.Env, cfg.getenv, cfg.appendFile)
}
