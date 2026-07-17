// Package model holds the leaf unified event types and portable hook result
// interfaces used by agnostic host adapters. Host packages under
// sdk/agnostic/internal/{claude,cursor,copilot} depend on model; they must not
// import sdk/agnostic (import cycle).
package model
