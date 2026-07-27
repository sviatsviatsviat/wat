package event

// AgentIdentity holds Copilot agent name fields shared by agent lifecycle events.
type AgentIdentity struct {
	// AgentName is the agent name.
	AgentName string `json:"agent_name"`
	// AgentDisplayName is the display name.
	AgentDisplayName string `json:"agent_display_name"`
}

// Name returns the agent name.
func (a AgentIdentity) Name() string {
	return a.AgentName
}

// DisplayName returns the agent display name.
func (a AgentIdentity) DisplayName() string {
	return a.AgentDisplayName
}
