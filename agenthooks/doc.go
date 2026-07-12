// Package agenthooks provides a unified hook event and decision model for
// Claude Code, GitHub Copilot, and Cursor. Handlers written against this
// package run unmodified under any supported agent; dialect codecs translate
// native stdin JSON into Event and Result back into agent-specific output.
package agenthooks
