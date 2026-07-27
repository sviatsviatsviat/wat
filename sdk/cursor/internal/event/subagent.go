package event

// SubagentFields holds Cursor subagent identity fields shared by subagent events.
type SubagentFields struct {
	// SubagentID is the subagent identifier when present.
	SubagentID string `json:"subagent_id"`
	// SubagentType is the subagent type. Official Hooks docs and hooks.json
	// matchers list camelCase values such as generalPurpose; live Cursor
	// payloads may use kebab-case such as general-purpose. Normalize before
	// comparing against matcher-style names.
	SubagentType string `json:"subagent_type"`
	// Task is the subagent task description.
	Task string `json:"task"`
}
