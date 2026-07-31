package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// Canonical update_todo / TodoWrite tool names accepted by [Input.AsUpdateTodo].
const (
	ToolUpdateTodo      = "update_todo"
	ToolTodoWriteClaude = "TodoWrite"
)

// UpdateTodoInput is the input schema for update_todo / TodoWrite.
// Fields follow the Claude TodoWrite shape (todos with content, status, activeForm).
type UpdateTodoInput struct {
	Todos []TodoItem `json:"todos"`
}

// TodoItem is one entry in UpdateTodoInput.
type TodoItem struct {
	Content    string `json:"content"`
	Status     string `json:"status"`
	ActiveForm string `json:"activeForm,omitempty"`
}

// AsUpdateTodo returns the update_todo tool input when this payload is for
// update_todo or TodoWrite.
func (in Input) AsUpdateTodo() (UpdateTodoInput, bool) {
	return hookkit.AsFold[UpdateTodoInput](in.Input, ToolUpdateTodo, ToolTodoWriteClaude)
}
