package shellparse

import (
	"runtime"
	"strings"

	"github.com/sviatsviatsviat/wat/internal/core"
)

const hostHintPowerShell = "powershell"

// hostHintFromGOOS returns hostHintPowerShell on Windows and an empty string elsewhere.
func hostHintFromGOOS() string {
	if runtime.GOOS == "windows" {
		return hostHintPowerShell
	}
	return ""
}

// detectDialect chooses "powershell" or "bash" using sniff markers, then host hint, then bash fallback.
func detectDialect(raw string, hostHint string) string {
	if hasPowerShellMarkers(raw) {
		return core.DialectPowerShell
	}
	if hasBashMarkers(raw) {
		return core.DialectBash
	}
	if hostHint == hostHintPowerShell {
		return core.DialectPowerShell
	}
	return core.DialectBash
}

func hasPowerShellMarkers(cmd string) bool {
	markers := []string{
		"Get-", "Set-", "New-", "Remove-", "Invoke-",
		"Select-Object", "Where-Object", "ForEach-Object",
		"$env:", "$PSVersionTable", "Write-Host", "Write-Output",
		"| Out-", "| Format-", "| ConvertTo-",
		"-ErrorAction", "-Verbose",
	}
	for _, m := range markers {
		if strings.Contains(cmd, m) {
			return true
		}
	}
	trimmed := strings.TrimSpace(cmd)
	return strings.HasPrefix(trimmed, "pwsh ") ||
		strings.HasPrefix(trimmed, "powershell ")
}

func hasBashMarkers(raw string) bool {
	markers := []string{"&&", "||", "/bin/", "/usr/", "sudo "}
	for _, m := range markers {
		if strings.Contains(raw, m) {
			return true
		}
	}
	return false
}
