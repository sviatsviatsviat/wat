# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Added

- Unified `agenthooks.Event` / `Kind` types and tool-name normalization for hook authors
- Unified `agenthooks.Result` / `Decision` types with `Merge` and `Unsupported` for hook handler responses
- `agenthooks.ParseDialect` and `agenthooks.Detect` for hook payload dialect sniffing
- `agenthooks.Codec`, `agenthooks.CodecFor`, and `agenthooks.ClaudeCodec` for Claude Code stdin/stdout translation
