package claude

import (
	"encoding/json"
	"fmt"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Encode renders a typed output struct as Claude Code stdout JSON.
// eventName is written into hookSpecificOutput.hookEventName.
func Encode(eventName string, out any, opts ...Option) ([]byte, error) {
	cfg := defaultRuntimeConfig()
	applyOptions(&cfg, opts...)
	if err := applyEnvSideEffect(eventName, out, cfg); err != nil {
		return nil, err
	}
	if eventName == "" {
		return nil, fmt.Errorf("claude: encode: empty event name")
	}
	out = hookkit.NormalizeOutput(out)
	if out == nil || isZeroOutput(out) {
		return nil, nil
	}
	if err := validateEncodePair(eventName, out); err != nil {
		return nil, err
	}

	top, hso, err := buildOutput(eventName, out)
	if err != nil {
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

func isZeroOutput(out any) bool {
	if z, ok := out.(interface{ isZero() bool }); ok {
		return z.isZero()
	}
	return hookkit.IsZeroOutput(out)
}

// IsZeroOutput reports whether out is an empty hook response.
func IsZeroOutput(out any) bool { return isZeroOutput(out) }

func validateEncodePair(eventName string, out any) error {
	allowed, ok := allowedEventsForOutput(out)
	if !ok {
		return fmt.Errorf("claude: encode: unsupported output type %T", out)
	}
	return hookkit.ValidateEncodePair("claude", eventName, out, allowed, nil)
}

func allowedEventsForOutput(out any) ([]string, bool) {
	switch out.(type) {
	case PreToolUseOutput:
		return []string{EventPreToolUse}, true
	case PermissionRequestOutput:
		return []string{EventPermissionRequest}, true
	case PostToolUseOutput:
		return []string{EventPostToolUse, EventPostToolUseFailure}, true
	case UserPromptSubmitOutput:
		return []string{EventUserPromptSubmit}, true
	case StopOutput:
		return []string{EventStop, EventSubagentStop}, true
	case SessionStartOutput:
		return []string{EventSessionStart}, true
	case MessageDisplayOutput:
		return []string{EventMessageDisplay}, true
	case PermissionDeniedOutput:
		return []string{EventPermissionDenied}, true
	case ElicitationOutput:
		return []string{EventElicitation}, true
	case WorktreeCreateOutput:
		return []string{EventWorktreeCreate}, true
	case CommonOutput:
		return []string{
			EventSetup,
			EventSessionEnd,
			EventUserPromptExpansion,
			EventPostToolBatch,
			EventSubagentStart,
			EventTaskCreated,
			EventTaskCompleted,
			EventStopFailure,
			EventTeammateIdle,
			EventNotification,
			EventInstructionsLoaded,
			EventConfigChange,
			EventCwdChanged,
			EventFileChanged,
			EventWorktreeRemove,
			EventPreCompact,
			EventPostCompact,
			EventElicitationResult,
		}, true
	default:
		return nil, false
	}
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
		return nil, nil, fmt.Errorf("claude: encode: unsupported output type %T", out)
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
