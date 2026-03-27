# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Initial **wat** CLI: first program argument selects the hook host (e.g. `cursor`), the second the wat subcommand (e.g. `run`), remaining arguments passed to that subcommand (e.g. `wat cursor run …`). **Cursor** is the supported host today; the same layout can extend to other hosts or subcommands (e.g. `guard`).
- **`run` subcommand** in package `run` (`run.NewRunCommand`): read hook JSON from stdin; optional `-f` / `--file-pattern` (Go regexp); when bindings include `__FILE_PATH__`, skip the subprocess if the cleaned path does not match; templated subprocess command with allowed placeholders; cross-platform execution with propagated child exit codes; hook stdout `{}` for Cursor, child stderr forwarded, child stdout discarded.
- Template variable substitution for the eight Cursor common envelope JSON fields (`conversation_id`, `generation_id`, `model`, `hook_event_name`, `cursor_version`, `workspace_roots`, `user_email`, `transcript_path`); event handling keyed on `hook_event_name`.
- Cursor hook types supported: `afterAgentResponse`, `afterAgentThought`, `afterFileEdit`, `afterMCPExecution`, `afterShellExecution`, `afterTabFileEdit`, `sessionEnd`.
- Unit and integration tests for CLI, event extraction, templating, and execution.
- GitHub Actions CI workflow for test, vet, and build checks.
- README and changelog documentation for usage, templatable parameters, and release baseline.
