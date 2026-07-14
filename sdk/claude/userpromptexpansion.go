package claude

// UserPromptExpansion is the UserPromptExpansion hook event.
type UserPromptExpansion struct {
	Envelope
	// ExpansionType is the expansion kind (slash_command, mcp_prompt).
	ExpansionType string `json:"expansion_type"`
	// CommandName is the slash command name.
	CommandName string `json:"command_name"`
	// CommandArgs is the slash command arguments.
	CommandArgs string `json:"command_args"`
	// CommandSource is the command source.
	CommandSource string `json:"command_source"`
	// Prompt is the expanded prompt text.
	Prompt string `json:"prompt"`
}

// EventName returns the hook event name.
func (UserPromptExpansion) EventName() string { return EventUserPromptExpansion }

func init() {
	registerDecoder(EventUserPromptExpansion, decodeAs[UserPromptExpansion])
}
