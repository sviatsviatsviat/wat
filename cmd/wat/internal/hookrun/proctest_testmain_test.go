package hookrun_test

import (
	"os"
	"testing"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/proctest"
)

func TestMain(m *testing.M) {
	proctest.MaybeExit()
	os.Exit(m.Run())
}
