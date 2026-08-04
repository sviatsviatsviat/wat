// Package proctest provides cross-platform subprocess stubs for CLI tests.
//
// Import this package only from *_test.go files. Call [MaybeExit] from each
// consuming package's TestMain so a child process started with WAT_PROCTEST_MODE
// can stand in for Unix utilities such as echo and sleep.
package proctest
