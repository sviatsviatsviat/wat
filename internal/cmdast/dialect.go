package cmdast

import (
	"runtime"
	"strings"
)

// Shell dialect labels returned by [DetectDialect] and
// [github.com/sviatsviatsviat/wat/internal/cmdast/hookline.ParseCommandLine].
const (
	DialectBash       = "bash"
	DialectPowerShell = "powershell"
)

// HostHintPowerShell is the host hint that steers [DetectDialect] toward PowerShell when markers are absent.
const HostHintPowerShell = "powershell"

// HostHintFromGOOS returns [HostHintPowerShell] on Windows and an empty string elsewhere.
func HostHintFromGOOS() string {
	if runtime.GOOS == "windows" {
		return HostHintPowerShell
	}
	return ""
}

// DetectDialect chooses [DialectPowerShell] or [DialectBash] using sniff markers,
// then host hint, then bash fallback. Use [HostHintFromGOOS] for the hook host hint on Windows.
func DetectDialect(raw string, hostHint string) string {
	if hasPowerShellMarkers(raw) {
		return DialectPowerShell
	}
	if hasBashMarkers(raw) {
		return DialectBash
	}
	if hostHint == HostHintPowerShell {
		return DialectPowerShell
	}
	return DialectBash
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
