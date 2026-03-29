# Lists agent transcript .jsonl paths that are missing from the incremental index or newer on disk
# than the indexed mtimeMs. Prints one path per line to stdout; writes count=N to stderr.
#
# Usage (from repo root); agent supplies the Cursor agent-transcripts directory:
#   powershell -File scripts/continual-learning/delta.ps1 -TranscriptsRoot "...\agent-transcripts"
#
# Exit: 0 on success; 1 if index unreadable.

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$TranscriptsRoot,
    [string]$IndexPath = ""
)

Set-StrictMode -Version Latest
$ErrorActionPreference = 'Stop'

. "$PSScriptRoot\common.ps1"

$TranscriptsRoot = [System.IO.Path]::GetFullPath($TranscriptsRoot)
if (-not (Test-Path -LiteralPath $TranscriptsRoot)) {
    Write-Error "Transcripts root directory not found: $TranscriptsRoot"
    exit 1
}
if ([string]::IsNullOrEmpty($IndexPath)) {
    $IndexPath = Join-Path (Get-WatRepoRoot) '.cursor/hooks/state/continual-learning-index.json'
}

try {
    $index = Read-ContinualLearningIndex -IndexPath $IndexPath
} catch {
    Write-Error "Cannot read index: $IndexPath - $_"
    exit 1
}

$indexed = $index.transcripts
$files = Get-AllJsonlFiles -TranscriptsRoot $TranscriptsRoot
$queue = [System.Collections.Generic.List[string]]::new()

foreach ($file in $files) {
    $p = Normalize-TranscriptPath -Path $file.FullName
    $diskMs = Get-FileMtimeMs -File $file
    $prev = $null
    if ($indexed.ContainsKey($p)) {
        $prev = $indexed[$p]
    }
    if ($null -eq $prev) {
        [void]$queue.Add($p)
        continue
    }
    $prevMs = $prev.mtimeMs
    if ($null -eq $prevMs) { [void]$queue.Add($p); continue }
    if ($diskMs -gt [long]$prevMs) {
        [void]$queue.Add($p)
    }
}

$count = $queue.Count
[Console]::Error.WriteLine(('count=' + $count))

foreach ($path in $queue) {
    Write-Output $path
}

exit 0
