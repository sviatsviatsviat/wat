# Shared helpers for continual-learning scripts (delta, sync-index, validate-jsonl).
# Repo root = parent of scripts/continual-learning.

function Get-WatRepoRoot {
    $twoUp = Join-Path (Join-Path $PSScriptRoot '..') '..'
    return (Resolve-Path $twoUp).Path
}

function Normalize-TranscriptPath {
    param([Parameter(Mandatory)][string]$Path)
    return [System.IO.Path]::GetFullPath($Path)
}

function Get-FileMtimeMs {
    param([Parameter(Mandatory)][System.IO.FileInfo]$File)
    # Avoid DateTimeOffset.SpecifyKind (not available on older runtimes); ticks are UTC wall time.
    $dto = [DateTimeOffset]::new($File.LastWriteTimeUtc.Ticks, [TimeSpan]::Zero)
    return $dto.ToUnixTimeMilliseconds()
}

function Read-ContinualLearningIndex {
    param([Parameter(Mandatory)][string]$IndexPath)
    if (-not (Test-Path -LiteralPath $IndexPath)) {
        return @{ version = 1; transcripts = @{} }
    }
    $raw = Get-Content -LiteralPath $IndexPath -Raw -Encoding UTF8
    if ([string]::IsNullOrWhiteSpace($raw)) {
        return @{ version = 1; transcripts = @{} }
    }
    $obj = $raw | ConvertFrom-Json
    $map = @{}
    if ($null -ne $obj.transcripts) {
        if ($obj.transcripts -is [hashtable]) {
            foreach ($key in $obj.transcripts.Keys) {
                $map[$key] = $obj.transcripts[$key]
            }
        } else {
            $obj.transcripts.PSObject.Properties | ForEach-Object {
                $map[$_.Name] = $_.Value
            }
        }
    }
    return @{ version = 1; transcripts = $map }
}

function Get-AllJsonlFiles {
    param([Parameter(Mandatory)][string]$TranscriptsRoot)
    if (-not (Test-Path -LiteralPath $TranscriptsRoot)) {
        return @()
    }
    return @(Get-ChildItem -LiteralPath $TranscriptsRoot -Filter '*.jsonl' -Recurse -File -ErrorAction Stop)
}
