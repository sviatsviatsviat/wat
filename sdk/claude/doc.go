// Package claude is the Claude Code hook SDK. Hook authors register typed
// handlers with UseHooks, then pass the chain to run.Serve from
// github.com/sviatsviatsviat/wat/sdk/run.
//
// Handlers build responses with hook-scoped *Results builders (and With* for
// advanced fields). A nil return is a silent no-op. Hook stdout is encoded
// internally from Output (sealed; only this package implements it), usually as
// JSON. WorktreeCreate is an exception: command hooks print a plain path.
//
// Blocking policy (see docs/agents/claude.md):
//
//   - Most denies use JSON decision fields (or permissionDecision / action) with
//     SuccessExit (0). Claude only applies JSON on exit 0.
//   - TeammateIdle, TaskCreated, and TaskCompleted Block builders use BlockExit
//     (2); the encoded reason is marked for stderr (Claude ignores stdout on
//     exit 2). WithContinue(false) stops the teammate entirely and is a
//     different control.
//   - Handler errors use HandlerErrorExit (1). Claude treats exit 1 as fail-open
//     for most events; WorktreeCreate fails on any non-zero exit.
//
// Tool input on PreToolUse and related events is typed as [Input] with AsBash,
// AsWrite, and related accessors (and Tool* name constants) in this package.
//
// PermissionRequest may carry permission_suggestions; allow responses can
// apply matching updatedPermissions via WithUpdatedPermissions using
// [PermissionUpdate] entries.
//
// See ExampleUseHooks in example_test.go for a minimal handler.
package claude
