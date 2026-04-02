$ErrorActionPreference = 'Stop'
# Command text is passed on standard input (not via environment variables).
$line = [Console]::In.ReadToEnd()
if ($null -eq $line) { $line = '' }

function ConvertSpan([System.Management.Automation.Language.IScriptExtent] $extent) {
    @{
        start = @{
            offset = [uint32]$extent.StartOffset
            line   = [uint32]$extent.StartLineNumber
            col    = [uint32]$extent.StartColumnNumber
        }
        end = @{
            offset = [uint32]$extent.EndOffset
            line   = [uint32]$extent.EndLineNumber
            col    = [uint32]$extent.EndColumnNumber
        }
    }
}

function ConvertVarRef([System.Management.Automation.Language.VariableExpressionAst] $v) {
    $scope = ''
    if ($null -ne $v.VariablePath -and $null -ne $v.VariablePath.DriveName) {
        $scope = [string]$v.VariablePath.DriveName
    }
    @{
        span  = (ConvertSpan $v.Extent)
        name  = [string]$v.VariablePath.UserPath
        scope = $scope
    }
}

function ConvertExpressionToArg($expr) {
    if ($null -eq $expr) { return $null }
    $tn = $expr.GetType().Name
    switch ($tn) {
        'StringConstantExpressionAst' {
            @{
                span    = (ConvertSpan $expr.Extent)
                literal = [string]$expr.Value
            }
        }
        'ConstantExpressionAst' {
            @{
                span    = (ConvertSpan $expr.Extent)
                literal = [string]$expr.Extent.Text
            }
        }
        'VariableExpressionAst' {
            @{
                span     = (ConvertSpan $expr.Extent)
                literal  = $expr.Extent.Text
                expanded = $true
                vars     = @( (ConvertVarRef $expr) )
            }
        }
        'ExpandableStringExpressionAst' {
            $vars = @()
            foreach ($n in $expr.NestedExpressions) {
                if ($n -is [System.Management.Automation.Language.VariableExpressionAst]) {
                    $vars += (ConvertVarRef $n)
                }
            }
            @{
                span     = (ConvertSpan $expr.Extent)
                literal  = [string]$expr.Value
                expanded = $true
                vars     = $vars
            }
        }
        'ScriptBlockExpressionAst' {
            @{
                span     = (ConvertSpan $expr.Extent)
                literal  = $expr.Extent.Text
                expanded = $true
            }
        }
        'SubExpressionAst' {
            @{
                span     = (ConvertSpan $expr.Extent)
                literal  = $expr.Extent.Text
                expanded = $true
            }
        }
        default {
            @{
                span     = (ConvertSpan $expr.Extent)
                literal  = $expr.Extent.Text
                expanded = $true
            }
        }
    }
}

function ConvertCommandAst([System.Management.Automation.Language.CommandAst] $node) {
    $elems = @($node.CommandElements)
    $name = ''
    if ($elems.Count -gt 0) {
        $name = $elems[0].Extent.Text.Trim()
    }
    $args = [System.Collections.Generic.List[object]]::new()
    $i = 1
    while ($i -lt $elems.Count) {
        $e = $elems[$i]
        if ($e -is [System.Management.Automation.Language.CommandParameterAst]) {
            $arg = @{
                span    = (ConvertSpan $e.Extent)
                literal = $e.Extent.Text
                flag    = @{
                    name   = [string]$e.ParameterName
                    dashes = 1
                }
            }
            if ($null -ne $e.Argument) {
                $arg.flag.value = (ConvertExpressionToArg $e.Argument)
            }
            $args.Add($arg)
        }
        else {
            $args.Add((ConvertExpressionToArg $e))
        }
        $i++
    }
    $redirs = @()
    foreach ($r in $node.Redirections) {
        $target = ''
        if ($r -is [System.Management.Automation.Language.FileRedirectionAst]) {
            if ($null -ne $r.Location) { $target = $r.Location.Extent.Text }
        }
        $opText = $r.Extent.Text
        if ($opText -match '^(\S+)') { $opText = $Matches[1] }
        $redirs += @{
            span     = (ConvertSpan $r.Extent)
            operator = [string]$opText
            target   = [string]$target
        }
    }
    @{
        name        = [string]$name
        args        = @($args)
        redirects   = $redirs
        assignments = @()
        background  = $false
        negated     = $false
    }
}

function PipelineElementToStatement($el) {
    if ($el -is [System.Management.Automation.Language.CommandAst]) {
        $cmd = ConvertCommandAst $el
        return @{
            kind    = 'Command'
            span    = (ConvertSpan $el.Extent)
            command = $cmd
        }
    }
    if ($el -is [System.Management.Automation.Language.PipelineAst]) {
        return (ConvertPipelineAst $el)
    }
    return @{
        kind     = 'Compound'
        span     = (ConvertSpan $el.Extent)
        compound = @{
            compound_kind = 'Other'
            raw             = $el.Extent.Text
        }
    }
}

function ConvertPipelineAst([System.Management.Automation.Language.PipelineAst] $node) {
    $els = @($node.PipelineElements)
    if ($els.Count -eq 1) {
        return (PipelineElementToStatement $els[0])
    }
    $stages = @()
    foreach ($el in $els) {
        $stages += (PipelineElementToStatement $el)
    }
    @{
        kind     = 'Pipeline'
        span     = (ConvertSpan $node.Extent)
        pipeline = @{
            stages     = $stages
            background = $false
        }
    }
}

function CompoundKindFromStatement($node) {
    $tn = $node.GetType().Name
    switch -Regex ($tn) {
        '^IfStatementAst$' { return 'If' }
        '^WhileStatementAst$' { return 'While' }
        '^ForStatementAst$' { return 'For' }
        '^ForEachStatementAst$' { return 'ForEach' }
        '^SwitchStatementAst$' { return 'Switch' }
        '^TryStatementAst$' { return 'Try' }
        '^FunctionDefinitionAst$' { return 'Function' }
        default { return 'Other' }
    }
}

function ConvertStatementAst($node) {
    if ($null -eq $node) { return $null }
    $tn = $node.GetType().Name
    switch ($tn) {
        'PipelineAst' { return (ConvertPipelineAst $node) }
        'CommandAst' {
            $cmd = ConvertCommandAst $node
            return @{
                kind    = 'Command'
                span    = (ConvertSpan $node.Extent)
                command = $cmd
            }
        }
        'AssignmentStatementAst' {
            return @{
                kind     = 'Compound'
                span     = (ConvertSpan $node.Extent)
                compound = @{
                    compound_kind = 'Other'
                    raw             = $node.Extent.Text
                }
            }
        }
        'PipelineChainAst' {
            $p = [System.Management.Automation.Language.PipelineChainAst]$node
            $pipes = @($p.Pipelines)
            $ops = @($p.OperatorTokens)
            if ($pipes.Count -eq 0) {
                return @{
                    kind     = 'Compound'
                    span     = (ConvertSpan $node.Extent)
                    compound = @{ compound_kind = 'Other'; raw = $node.Extent.Text }
                }
            }
            $acc = PipelineElementToStatement $pipes[0]
            for ($i = 0; $i -lt $ops.Count; $i++) {
                $tok = $ops[$i].Kind
                $opStr = '&&'
                if ($tok -eq [System.Management.Automation.Language.TokenKind]::OrOr) { $opStr = '||' }
                $right = PipelineElementToStatement $pipes[$i + 1]
                $acc = @{
                    kind  = 'Chain'
                    span  = @{
                        start = $acc.span.start
                        end   = $right.span.end
                    }
                    chain = @{
                        operator = $opStr
                        left     = $acc
                        right    = $right
                    }
                }
            }
            return $acc
        }
        default {
            $ck = CompoundKindFromStatement $node
            return @{
                kind     = 'Compound'
                span     = (ConvertSpan $node.Extent)
                compound = @{
                    compound_kind = [string]$ck
                    raw             = $node.Extent.Text
                }
            }
        }
    }
}

$tokens = $null
$errors = $null
$ast = [System.Management.Automation.Language.Parser]::ParseInput($line, [ref]$tokens, [ref]$errors)

$stmts = @()
if ($null -ne $ast.EndBlock) {
    foreach ($s in $ast.EndBlock.Statements) {
        $stmts += (ConvertStatementAst $s)
    }
}

$result = @{
    lang         = 'powershell'
    raw          = [string]$line
    statements   = $stmts
    parse_error  = ''
    parse_partial = $false
}

if ($errors.Count -gt 0) {
    $result.parse_partial = $true
    $result.parse_error = ($errors | ForEach-Object { $_.Message }) -join '; '
}

($result | ConvertTo-Json -Depth 80 -Compress)
