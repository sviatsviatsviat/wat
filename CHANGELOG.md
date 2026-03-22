# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Initial `wat run` command with stdin hook-event parsing; subprocess argv after the subcommand; **Cursor** as default host with optional `--host` / `-H` after the subcommand (e.g. `wat run -H cursor …`); same pattern for other commands such as `guard`.
- Template variable substitution for the eight Cursor common envelope JSON fields (`conversation_id`, `generation_id`, `model`, `hook_event_name`, `cursor_version`, `workspace_roots`, `user_email`, `transcript_path`); event handling keyed on `hook_event_name`.
- Cursor hook types supported: `afterAgentResponse`, `afterAgentThought`, `afterFileEdit`, `afterMCPExecution`, `afterShellExecution`, `afterTabFileEdit`, `sessionEnd`.
- Cross-platform command execution with propagated exit codes; hook stdout emits `{}` for Cursor JSON parsing while the child’s stderr is forwarded to wat’s stderr and the child’s stdout is discarded.
- Unit and integration tests for CLI, event extraction, templating, and execution.
- GitHub Actions CI workflow for test, vet, and build checks.
- README and changelog documentation for usage, templatable parameters, and release baseline.
