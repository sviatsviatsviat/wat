package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolAskUserQuestion is the ask-user-question tool.
const ToolAskUserQuestion = "AskUserQuestion"

// AskUserQuestionInput is the input schema for the AskUserQuestion tool.
type AskUserQuestionInput struct {
	Questions []Question        `json:"questions"`
	Answers   map[string]string `json:"answers,omitempty"`
}

// Question is one question in AskUserQuestionInput.
type Question struct {
	Question    string   `json:"question"`
	Header      string   `json:"header,omitempty"`
	Options     []Option `json:"options,omitempty"`
	MultiSelect bool     `json:"multiSelect,omitempty"`
}

// Option is one selectable option in a Question.
type Option struct {
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// AsAskUserQuestion returns the AskUserQuestion tool input when this payload is for AskUserQuestion.
func (in Input) AsAskUserQuestion() (AskUserQuestionInput, bool) {
	return hookkit.As[AskUserQuestionInput](in.Input, ToolAskUserQuestion)
}
