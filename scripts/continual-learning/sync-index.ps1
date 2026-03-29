# Rebuilds .cursor/hooks/state/continual-learning-index.json from all *.jsonl under the
# transcripts tree: sets mtimeMs from disk, lastProcessedAt to UTC now (or -ProcessedAt),
# drops index entries for deleted files.
#
# Usage; agent supplies the Cursor agent-transcripts directory:
#   powershell -File scripts/continual-learning/sync-index.ps1 -TranscriptsRoot "...\agent-transcripts"
#
# Exit: 0 on success; non-zero on I/O or JSON errors.

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$TranscriptsRoot,
    [string]$IndexPath = "",
    [string]$ProcessedAt = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

. "$PSScriptRoot\common.ps1"

$TranscriptsRoot = [System.IO.Path]::GetFullPath($TranscriptsRoot)
if ([string]::IsNullOrEmpty($IndexPath)) {
    $IndexPath = Join-Path (Get-WatRepoRoot) '.cursor\hooks\state\continual-learning-index.json'
}
if ([string]::IsNullOrEmpty($ProcessedAt)) {
    $ProcessedAt = [DateTimeOffset]::UtcNow.ToString('yyyy-MM-ddTHH:mm:ss.fffZ')
}

$files = Get-AllJsonlFiles -TranscriptsRoot $TranscriptsRoot
$transcripts = @{}
foreach ($file in $files) {
    $p = Normalize-TranscriptPath -Path $file.FullName
    $ms = Get-FileMtimeMs -File $file
    $transcripts[$p] = @{
        mtimeMs           = $ms
        lastProcessedAt = $ProcessedAt
    }
}

$root = [ordered]@{
    version     = 1
    transcripts = $transcripts
}

$json = ($root | ConvertTo-Json -Depth 10)
$dir = Split-Path -Parent $IndexPath
if (-not (Test-Path -LiteralPath $dir)) {
    New-Item -ItemType Directory -Path $dir -Force | Out-Null
}

$utf8NoBom = New-Object System.Text.UTF8Encoding $false
[System.IO.File]::WriteAllText($IndexPath, $json + "`n", $utf8NoBom)

Write-Output "Wrote $IndexPath ($($transcripts.Count) transcripts)."
exit 0
