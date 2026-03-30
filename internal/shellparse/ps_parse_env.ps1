$ErrorActionPreference = 'Stop'
$line = $env:WAT_CMD
if ($null -eq $line) { Write-Output '[]'; exit 0 }
$tokens = $null; $errors = $null
$ast = [System.Management.Automation.Language.Parser]::ParseInput($line, [ref]$tokens, [ref]$errors)
if ($errors.Count -gt 0) { Write-Output '[]'; exit 0 }
$pipeline = $ast.FindAll({ param($n) $n -is [System.Management.Automation.Language.CommandAst] }, $true)
$pipeLen = 1
$pipeAsts = $ast.FindAll({ param($n) $n -is [System.Management.Automation.Language.PipelineAst] }, $false)
if ($pipeAsts.Count -gt 0 -and $pipeAsts[0].PipelineElements.Count -gt 0) {
  $pipeLen = $pipeAsts[0].PipelineElements.Count
}
$result = [System.Collections.Generic.List[hashtable]]::new()
$idx = 0
foreach ($cmd in $pipeline) {
  $name = $cmd.GetCommandName()
  $elements = @($cmd.CommandElements | Select-Object -Skip 1)
  $flags = @{}
  $switches = [System.Collections.Generic.List[string]]::new()
  $posArgs = [System.Collections.Generic.List[string]]::new()
  for ($i = 0; $i -lt $elements.Count; $i++) {
    $el = $elements[$i]
    if ($el -is [System.Management.Automation.Language.CommandParameterAst]) {
      if ($null -ne $el.Argument) {
        $flags[$el.ParameterName] = $el.Argument.Extent.Text
      } elseif ($i + 1 -lt $elements.Count -and $elements[$i + 1] -isnot [System.Management.Automation.Language.CommandParameterAst]) {
        $flags[$el.ParameterName] = $elements[$i + 1].Extent.Text
        $i++
      } else {
        $switches.Add($el.ParameterName)
      }
    } else {
      $posArgs.Add($el.Extent.Text)
    }
  }
  $result.Add(@{
    Name = $name
    Args = @($posArgs)
    Flags = $flags
    Switches = @($switches)
    PipeIndex = $idx
    PipeLength = $pipeLen
  })
  $idx++
}
($result | ConvertTo-Json -Depth 8 -Compress)
