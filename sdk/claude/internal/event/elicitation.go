package event

// ElicitationFields holds Claude elicitation wire fields shared by
// Elicitation and ElicitationResult events.
type ElicitationFields struct {
	// MCPServerName is the MCP server involved in the elicitation.
	MCPServerName string `json:"mcp_server_name"`
	// Mode is the elicitation mode when provided ("form" or "url").
	Mode string `json:"mode"`
	// ElicitationID is the elicitation request identifier when provided.
	ElicitationID string `json:"elicitation_id"`
}
