package hookkit

import "encoding/json"

// ExtractShellCommand reads the command string from a shell tool input object.
func ExtractShellCommand(input json.RawMessage) string {
	if len(input) == 0 {
		return ""
	}
	var args struct {
		Command string `json:"command"`
	}
	if json.Unmarshal(input, &args) != nil {
		return ""
	}
	return args.Command
}
