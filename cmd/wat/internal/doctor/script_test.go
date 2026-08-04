package doctor

import (
	"errors"
	"testing"
	"time"

	"github.com/sviatsviatsviat/wat/cmd/wat/internal/proctest"
)

func TestCombinedOutputWithTimeout_success(t *testing.T) {
	cmd := proctest.Println("ok")
	out, err := combinedOutputWithTimeout(cmd, time.Second)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "ok\n" && string(out) != "ok\r\n" {
		t.Fatalf("out = %q", out)
	}
}

func TestCombinedOutputWithTimeout_timeout(t *testing.T) {
	cmd := proctest.Sleep()
	_, err := combinedOutputWithTimeout(cmd, 50*time.Millisecond)
	if !errors.Is(err, errCommandTimeout) {
		t.Fatalf("err = %v, want timeout", err)
	}
}
