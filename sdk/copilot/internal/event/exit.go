package event

// HandlerErrorExit is exit code 1 for handler errors. Copilot command hooks
// fail-closed on non-zero exits other than 2 (including PreToolUse).
const HandlerErrorExit = 1

// WarnExit is exit code 2. Copilot treats it as a warning by default; for
// PermissionRequest it means deny, and for PostToolUseFailure it carries
// additional_context in stdout.
const WarnExit = 2
