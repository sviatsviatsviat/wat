package claude

// HandlerErrorExit is exit code 1 for mux processing failures (read/decode/handler/encode/write errors).
const HandlerErrorExit = 1

// FailBlockExit is exit code 2 for handler errors when WithFailPolicy(FailBlock) is active.
const FailBlockExit = 2
