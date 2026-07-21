package claude

// SuccessExit is exit code 0 for successful encode and no-op responses.
const SuccessExit = 0

// HandlerErrorExit is exit code 1 for mux processing failures (read/decode/handler/encode/write errors).
const HandlerErrorExit = 1
