package event

// PermissionUpdateType is the discriminator for a permission update entry.
type PermissionUpdateType string

const (
	// PermissionUpdateAddRules adds permission rules.
	PermissionUpdateAddRules PermissionUpdateType = "addRules"
	// PermissionUpdateReplaceRules replaces rules of a given behavior at a destination.
	PermissionUpdateReplaceRules PermissionUpdateType = "replaceRules"
	// PermissionUpdateRemoveRules removes matching rules of a given behavior.
	PermissionUpdateRemoveRules PermissionUpdateType = "removeRules"
	// PermissionUpdateSetMode changes the permission mode.
	PermissionUpdateSetMode PermissionUpdateType = "setMode"
	// PermissionUpdateAddDirectories adds working directories.
	PermissionUpdateAddDirectories PermissionUpdateType = "addDirectories"
	// PermissionUpdateRemoveDirectories removes working directories.
	PermissionUpdateRemoveDirectories PermissionUpdateType = "removeDirectories"
)

// PermissionDestination controls where a permission update is written.
type PermissionDestination string

const (
	// PermissionDestinationSession keeps the update in memory for the session.
	PermissionDestinationSession PermissionDestination = "session"
	// PermissionDestinationLocalSettings writes to .claude/settings.local.json.
	PermissionDestinationLocalSettings PermissionDestination = "localSettings"
	// PermissionDestinationProjectSettings writes to .claude/settings.json.
	PermissionDestinationProjectSettings PermissionDestination = "projectSettings"
	// PermissionDestinationUserSettings writes to ~/.claude/settings.json.
	PermissionDestinationUserSettings PermissionDestination = "userSettings"
)

// PermissionRule is one tool rule inside an add/replace/removeRules update.
type PermissionRule struct {
	// ToolName is the tool the rule matches.
	ToolName string `json:"toolName"`
	// RuleContent is an optional tool-specific pattern; omit to match the whole tool.
	RuleContent string `json:"ruleContent,omitempty"`
}

// PermissionUpdate is one Claude permission update / suggestion entry.
//
// PermissionRequest input uses these entries as permission_suggestions.
// Allow responses may echo them as updatedPermissions.
type PermissionUpdate struct {
	// Type selects which fields apply (addRules, setMode, addDirectories, …).
	Type PermissionUpdateType `json:"type"`
	// Rules is the rule list for addRules, replaceRules, and removeRules.
	Rules []PermissionRule `json:"rules,omitempty"`
	// Behavior is allow, deny, or ask for rule updates.
	Behavior PermissionDecision `json:"behavior,omitempty"`
	// Destination controls whether the change is session-only or persisted.
	Destination PermissionDestination `json:"destination,omitempty"`
	// Mode is the permission mode for setMode entries.
	Mode PermissionMode `json:"mode,omitempty"`
	// Directories is the path list for addDirectories and removeDirectories.
	Directories []string `json:"directories,omitempty"`
}

// ClonePermissionUpdates returns a deep copy of updates.
// Nil stays nil. Nested Rules and Directories slices are cloned.
func ClonePermissionUpdates(updates []PermissionUpdate) []PermissionUpdate {
	if updates == nil {
		return nil
	}
	out := make([]PermissionUpdate, len(updates))
	for i, u := range updates {
		out[i] = u
		if u.Rules != nil {
			out[i].Rules = append([]PermissionRule(nil), u.Rules...)
		}
		if u.Directories != nil {
			out[i].Directories = append([]string(nil), u.Directories...)
		}
	}
	return out
}
