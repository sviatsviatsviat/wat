package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical ask-user tool names accepted by [Input.AsAskUser].
const (
	ToolAskUser               = "ask_user"
	ToolAskUserQuestionClaude = "AskUserQuestion"
)

// AskUserInput is the input schema for the ask_user / AskUserQuestion tool.
type AskUserInput struct {
	Questions []AskUserQuestion `json:"questions"`
	Answers   map[string]string `json:"answers,omitempty"`
}

// AskUserQuestion is one question in AskUserInput.
type AskUserQuestion struct {
	Question    string          `json:"question"`
	Header      string          `json:"header,omitempty"`
	Options     []AskUserOption `json:"options,omitempty"`
	MultiSelect bool            `json:"multiSelect,omitempty"`
}

// AskUserOption is one selectable option in an AskUserQuestion.
type AskUserOption struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// AsAskUser returns the ask_user tool input when this payload is for ask_user or AskUserQuestion.
func (in Input) AsAskUser() (AskUserInput, bool) {
	return hookkit.AsFold[AskUserInput](in.Input, ToolAskUser, ToolAskUserQuestionClaude)
}
