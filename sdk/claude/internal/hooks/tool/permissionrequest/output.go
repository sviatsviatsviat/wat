package permissionrequest

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
	"github.com/sviatsviatsviat/wat/sdk/claude/internal/event"
)

// Output is the response for this hook event.
// Construct via Results builders and With* methods.
// A nil value is a no-op.
type Output interface {
	hookkit.Output
	isOutput()
	// WithUpdatedInput replaces tool arguments when set.
	WithUpdatedInput(input map[string]any) Output
	// WithUpdatedPermissions applies permission update entries on allow.
	// Entries match the permission_suggestions shape; a common pattern is to
	// echo a suggestion from the event. The slice is cloned.
	WithUpdatedPermissions(updates []event.PermissionUpdate) Output
	// WithInterrupt stops the session when true.
	WithInterrupt(v bool) Output
	// WithAdditionalContext injects model context.
	WithAdditionalContext(text string) Output
	// WithContinue sets whether Claude should continue the session.
	WithContinue(v bool) Output
	// WithStopReason explains why the session was stopped.
	WithStopReason(reason string) Output
	// WithSuppressOutput suppresses hook stdout when true.
	WithSuppressOutput(v bool) Output
	// WithSystemMessage sets a user-visible system message.
	WithSystemMessage(msg string) Output
	// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
	WithTerminalSequence(seq string) Output
}

type output struct {
	event.Common
	behavior           string
	updatedInput       map[string]any
	updatedPermissions []event.PermissionUpdate
	message            string
	interrupt          bool
	additionalContext  string
}

func (output) isOutput() {}

// IsZero reports whether this hook response is empty.
func (o output) IsZero() bool {
	return o.Common.IsZero() && o.behavior == "" && o.updatedInput == nil &&
		o.updatedPermissions == nil && o.message == "" && !o.interrupt &&
		o.additionalContext == ""
}

// WithUpdatedInput replaces tool arguments when set.
func (o output) WithUpdatedInput(input map[string]any) Output {
	o.updatedInput = input
	return o
}

// WithUpdatedPermissions applies permission update entries on allow.
func (o output) WithUpdatedPermissions(updates []event.PermissionUpdate) Output {
	o.updatedPermissions = event.ClonePermissionUpdates(updates)
	return o
}

// WithInterrupt stops the session when true.
func (o output) WithInterrupt(v bool) Output {
	o.interrupt = v
	return o
}

// WithAdditionalContext injects model context.
func (o output) WithAdditionalContext(text string) Output {
	o.additionalContext = text
	return o
}

// WithContinue sets whether Claude should continue the session.
func (o output) WithContinue(v bool) Output {
	o.Common = o.Common.WithContinue(v)
	return o
}

// WithStopReason explains why the session was stopped.
func (o output) WithStopReason(reason string) Output {
	o.Common = o.Common.WithStopReason(reason)
	return o
}

// WithSuppressOutput suppresses hook stdout when true.
func (o output) WithSuppressOutput(v bool) Output {
	o.Common = o.Common.WithSuppressOutput(v)
	return o
}

// WithSystemMessage sets a user-visible system message.
func (o output) WithSystemMessage(msg string) Output {
	o.Common = o.Common.WithSystemMessage(msg)
	return o
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (o output) WithTerminalSequence(seq string) Output {
	o.Common = o.Common.WithTerminalSequence(seq)
	return o
}

func (o output) encodeInto(top, hso map[string]any) {
	event.ApplyCommon(top, o.Common)
	if o.behavior != "" {
		dec := map[string]any{"behavior": o.behavior}
		if o.updatedInput != nil {
			dec["updatedInput"] = o.updatedInput
		}
		if o.behavior == "allow" && o.updatedPermissions != nil {
			dec["updatedPermissions"] = o.updatedPermissions
		}
		if o.message != "" {
			dec["message"] = o.message
		}
		if o.interrupt {
			dec["interrupt"] = true
		}
		hso["decision"] = dec
	}
	if o.additionalContext != "" {
		hso["additionalContext"] = o.additionalContext
	}
}

// Encode renders this output as Claude Code stdout JSON.
func (o output) Encode() ([]byte, int, error) {
	return event.MarshalHookOutput(event.PermissionRequest, o.encodeInto)
}

// Merge combines other into this PermissionRequest output.
func (o output) Merge(other hookkit.Output) (hookkit.Output, []string, error) {
	b, ok := other.(output)
	if !ok {
		return nil, nil, hookkit.ErrMergeType(o, other)
	}
	mergedCommon, warnings := o.Common.Merge(b.Common)
	behavior, message := hookkit.MergeRankedString(
		o.behavior, o.message,
		b.behavior, b.message,
		hookkit.PermissionRankString,
	)
	updatedInput, w := hookkit.TakeLastMap("updatedInput", o.updatedInput, b.updatedInput)
	if w != "" {
		warnings = append(warnings, w)
	}
	updatedPermissions, w := takeLastPermissionUpdates(o.updatedPermissions, b.updatedPermissions)
	if w != "" {
		warnings = append(warnings, w)
	}
	return output{
		Common:             mergedCommon,
		behavior:           behavior,
		updatedInput:       updatedInput,
		updatedPermissions: updatedPermissions,
		message:            message,
		interrupt:          o.interrupt || b.interrupt,
		additionalContext:  hookkit.JoinContextStrings(o.additionalContext, b.additionalContext),
	}, warnings, nil
}

// Stop reports whether remaining handlers should be skipped.
func (o output) Stop() bool {
	return o.Common.Stop() || o.behavior == "deny"
}

func takeLastPermissionUpdates(dst, src []event.PermissionUpdate) ([]event.PermissionUpdate, string) {
	if src == nil {
		return dst, ""
	}
	cloned := event.ClonePermissionUpdates(src)
	if dst != nil {
		return cloned, hookkit.OverwriteWarning("updatedPermissions")
	}
	return cloned, ""
}
