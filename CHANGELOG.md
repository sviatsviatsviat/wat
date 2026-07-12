# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Unified `agenthooks.Event` / `Kind` types and tool-name normalization for hook authors
- Unified `agenthooks.Result` / `Decision` types with `Merge` and `Unsupported` for hook handler responses
- `agenthooks.ParseDialect` and `agenthooks.Detect` for hook payload dialect sniffing
- `agenthooks.Codec` and `agenthooks.CodecFor` for dialect codec lookup
- `agenthooks.ClaudeCodec` for Claude Code stdin/stdout translation
- `agenthooks.CopilotCodec` for GitHub Copilot camelCase and VS Code hook payloads
- `agenthooks.CursorCodec` decodes and encodes Cursor hooks, folding dedicated shell/MCP/file events into unified pre/post tool kinds
- `agenthooks.Mux` with `On`, `OnAny`, `Serve`, and `Main` for hook handler registration and the hook process lifecycle
- Serve options `WithDialect`, `WithEvent`, and `WithGetenv` for explicit dialect, Copilot event hints, and testable environment lookup
- `agenthooks/portconfig.Parse` and `Emit` for Claude Code, GitHub Copilot, and Cursor hook configuration files, preserving unmappable entries and native handler fields for same-dialect round-trip
