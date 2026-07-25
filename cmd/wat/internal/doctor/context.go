package doctor

import "github.com/sviatsviatsviat/wat/sdk/run"

// Context carries resolved project state shared across doctor checks.
type Context struct {
	WatDir      string
	WatErr      error
	Manifest    run.Manifest
	ManifestErr error
}
