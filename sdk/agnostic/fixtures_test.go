package agnostic

const claudePreToolUse = `{
  "session_id": "abc123",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/home/user/proj",
  "permission_mode": "default",
  "hook_event_name": "PreToolUse",
  "tool_name": "Bash",
  "tool_use_id": "tu_1",
  "tool_input": {"command": "rm -rf /tmp/build", "description": "clean"}
}`

const copilotPreToolUse = `{
  "hook_event_name": "PreToolUse",
  "session_id": "s1",
  "timestamp": "2026-07-12T10:00:00Z",
  "cwd": "/w",
  "tool_name": "bash",
  "tool_input": {"command": "rm -rf /"}
}`

const cursorShell = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "some-model",
  "hook_event_name": "beforeShellExecution",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "user_email": null,
  "transcript_path": null,
  "command": "git push --force",
  "cwd": "/w",
  "sandbox": false
}`

const cursorPreToolUse = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "some-model",
  "hook_event_name": "preToolUse",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "tool_name": "Shell",
  "tool_input": {"command": "rm -rf /"},
  "tool_use_id": "tu1",
  "cwd": "/w"
}`

const cursorBeforeRead = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "some-model",
  "hook_event_name": "beforeReadFile",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "file_path": "secret.env",
  "content": "KEY=1",
  "cwd": "/w"
}`

const cursorAfterShell = `{
  "conversation_id": "c1",
  "generation_id": "g1",
  "model": "some-model",
  "hook_event_name": "afterShellExecution",
  "cursor_version": "1.7.2",
  "workspace_roots": ["/w"],
  "command": "echo hi",
  "output": "hi",
  "duration": 12,
  "cwd": "/w"
}`

const claudeStop = `{
  "session_id": "abc123",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/home/user/proj",
  "permission_mode": "default",
  "hook_event_name": "Stop",
  "stop_hook_active": false
}`

const claudeSessionStart = `{
  "session_id": "abc123",
  "transcript_path": "/tmp/t.jsonl",
  "cwd": "/home/user/proj",
  "permission_mode": "default",
  "hook_event_name": "SessionStart",
  "source": "startup"
}`
