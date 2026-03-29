# Reads each line of every *.jsonl under the transcripts root as JSON; prints errors as path:line: message.
# Empty lines are skipped. Exit 1 if any line failed; 0 if all lines parsed (or no files).
#
# Usage; agent supplies the Cursor agent-transcripts directory:
#   powershell -File scripts/continual-learning/validate-jsonl.ps1 -TranscriptsRoot "...\agent-transcripts"
#
# Exit: 0 if all lines valid; 1 if parse failed or root missing on disk.

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$TranscriptsRoot
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

. "$PSScriptRoot\common.ps1"

$TranscriptsRoot = [System.IO.Path]::GetFullPath($TranscriptsRoot)

if (-not (Test-Path -LiteralPath $TranscriptsRoot)) {
    Write-Error "Transcripts root not found: $TranscriptsRoot"
    exit 1
}

$failed = $false
$files = Get-AllJsonlFiles -TranscriptsRoot $TranscriptsRoot

foreach ($file in $files) {
    $lineNum = 0
    $reader = [System.IO.StreamReader]::new($file.FullName)
    try {
        while ($null -ne ($line = $reader.ReadLine())) {
            $lineNum++
            if ([string]::IsNullOrWhiteSpace($line)) { continue }
            try {
                $null = $line | ConvertFrom-Json
            } catch {
                Write-Output "$($file.FullName):${lineNum}: $_"
                $failed = $true
            }
        }
    } finally {
        $reader.Close()
    }
}

if ($failed) { exit 1 }
exit 0
