package event

import (
	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Common holds output fields shared across Claude Code hook responses.
type Common struct {
	Cont             *bool
	StopReason       string
	SuppressOutput   bool
	SystemMessage    string
	TerminalSequence string
}

// IsZero reports whether this hook response is empty.
func (c Common) IsZero() bool {
	return c.Cont == nil && c.StopReason == "" && !c.SuppressOutput &&
		c.SystemMessage == "" && c.TerminalSequence == ""
}

// WithContinue sets whether Claude should continue the session.
// Pass false to stop Claude entirely.
func (c Common) WithContinue(v bool) Common {
	c.Cont = &v
	return c
}

// WithStopReason explains why the session was stopped.
func (c Common) WithStopReason(reason string) Common {
	c.StopReason = reason
	return c
}

// WithSuppressOutput suppresses hook stdout when true.
func (c Common) WithSuppressOutput(v bool) Common {
	c.SuppressOutput = v
	return c
}

// WithSystemMessage sets a user-visible system message.
func (c Common) WithSystemMessage(msg string) Common {
	c.SystemMessage = msg
	return c
}

// WithTerminalSequence sets an OSC terminal sequence (allowlisted).
func (c Common) WithTerminalSequence(seq string) Common {
	c.TerminalSequence = seq
	return c
}

// Merge combines other into this shared fields set.
func (c Common) Merge(other Common) (Common, []string) {
	var warnings []string
	out := c

	// continue:false is sticky; otherwise last non-nil wins.
	if c.Cont != nil && !*c.Cont {
		out.Cont = c.Cont
		if other.Cont != nil && !*other.Cont {
			sr, w := hookkit.TakeLastString("stopReason", c.StopReason, other.StopReason)
			out.StopReason = sr
			if w != "" {
				warnings = append(warnings, w)
			}
		} else {
			out.StopReason = c.StopReason
		}
	} else if other.Cont != nil {
		if c.Cont != nil {
			warnings = append(warnings, hookkit.OverwriteWarning("continue"))
		}
		out.Cont = other.Cont
		if !*other.Cont {
			out.StopReason = other.StopReason
		} else {
			sr, w := hookkit.TakeLastString("stopReason", c.StopReason, other.StopReason)
			out.StopReason = sr
			if w != "" {
				warnings = append(warnings, w)
			}
		}
	} else {
		sr, w := hookkit.TakeLastString("stopReason", c.StopReason, other.StopReason)
		out.StopReason = sr
		if w != "" {
			warnings = append(warnings, w)
		}
	}

	out.SuppressOutput = c.SuppressOutput || other.SuppressOutput

	sm, w := hookkit.TakeLastString("systemMessage", c.SystemMessage, other.SystemMessage)
	out.SystemMessage = sm
	if w != "" {
		warnings = append(warnings, w)
	}
	ts, w := hookkit.TakeLastString("terminalSequence", c.TerminalSequence, other.TerminalSequence)
	out.TerminalSequence = ts
	if w != "" {
		warnings = append(warnings, w)
	}
	return out, warnings
}

// Stop reports whether remaining handlers should be skipped.
func (c Common) Stop() bool {
	return c.Cont != nil && !*c.Cont
}

// ApplyCommon writes shared response fields into top.
func ApplyCommon(top map[string]any, c Common) {
	if c.Cont != nil {
		top["continue"] = *c.Cont
		if !*c.Cont && c.StopReason != "" {
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
