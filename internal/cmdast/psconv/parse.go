package psconv

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/cmdast"
)

// serializeScript is the embedded PowerShell serializer; it reads the command line from stdin
// and prints one JSON object matching [cmdast.CommandLine].
//
//go:embed serialize_cmdast.ps1
var serializeScript string

// Parse runs PowerShell (**pwsh** preferred, else **powershell** on PATH) with the embedded script,
// passes raw on **stdin**, and unmarshals JSON into [cmdast.CommandLine].
// There is no fallback tokenizer for this AST path.
func Parse(raw string) (*cmdast.CommandLine, error) {
	exe, err := resolvePowerShell()
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(exe, "-NoProfile", "-NoLogo", "-NonInteractive", "-Command", serializeScript)
	cmd.Stdin = strings.NewReader(raw)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return nil, fmt.Errorf("powershell cmdast: %w: %s", err, msg)
		}
		return nil, fmt.Errorf("powershell cmdast: %w", err)
	}
	payload := strings.TrimSpace(string(out))
	if payload == "" {
		return nil, fmt.Errorf("powershell cmdast: empty output")
	}
	var cl cmdast.CommandLine
	if err := json.Unmarshal([]byte(payload), &cl); err != nil {
		return nil, fmt.Errorf("cmdast json: %w", err)
	}
	cl.Lang = cmdast.LangPowerShell
	if cl.Raw == "" {
		cl.Raw = raw
	}
	return &cl, nil
}

func resolvePowerShell() (string, error) {
	for _, name := range []string{"pwsh", "powershell"} {
		path, err := exec.LookPath(name)
		if err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("powershell cmdast: neither pwsh nor powershell on PATH")
}
