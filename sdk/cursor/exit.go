package cursor

// HandlerErrorExit is exit code 1. The runner should use this when a handler
// returns an error under Cursor's default fail-open policy.
const HandlerErrorExit = 1

// PermissionDenyExit is exit code 2. Cursor treats it as block/deny on permission-gating
// events, equivalent to returning permission:"deny".
const PermissionDenyExit = 2
