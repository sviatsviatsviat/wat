package tools

import "github.com/sviatsviatsviat/wat/internal/hookkit"

// ToolExitPlanMode is the exit-plan-mode tool.
const ToolExitPlanMode = "ExitPlanMode"

// ExitPlanModeInput is the input schema for the ExitPlanMode tool.
type ExitPlanModeInput struct {
	Plan         string `json:"plan"`
	PlanFilePath string `json:"plan_file_path,omitempty"`
}

// AsExitPlanMode returns the ExitPlanMode tool input when this payload is for ExitPlanMode.
func (in Input) AsExitPlanMode() (ExitPlanModeInput, bool) {
	return hookkit.As[ExitPlanModeInput](in.Input, ToolExitPlanMode)
}
