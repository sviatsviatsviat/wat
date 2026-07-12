// Package copilothook is the GitHub Copilot hook SDK. Hook authors register
// typed handlers; the SDK decodes stdin JSON (camelCase and VS Code compatible
// formats), dispatches by event, encodes stdout, and selects exit codes.
package copilothook
