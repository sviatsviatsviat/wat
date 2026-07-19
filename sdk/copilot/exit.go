package copilot

// PreToolErrorExit is the exit code when a PreToolUse handler returns an error.
// Copilot command hooks fail-closed on non-zero exits other than 2.
const PreToolErrorExit = 1

// HandlerErrorExit is exit code 1 for handler errors under fail-open policy.
const HandlerErrorExit = PreToolErrorExit

// WarnExit is exit code 2. Copilot treats it as a warning by default; for
// PermissionRequest it means deny, and for PostToolUseFailure it carries
// additional_context in stdout.
const WarnExit = 2
