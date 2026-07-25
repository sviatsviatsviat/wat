// Package dotwat tests this repository's own .wat/ hook project — the
// dogfooded hooks.go and go.mod committed under .wat/ — by driving it through
// the wat CLI and testdata fixtures, the same way an end user would.
//
// e2e scaffolds a fresh project from the wat init template and knows nothing
// about this repository's committed .wat/hooks.go; that file lives in a
// separate "wat-hooks" Go module (see .wat/go.mod) that "go test ./..." at
// the repo root does not compile, so this package is the only place its
// behavior is asserted in CI.
package dotwat
