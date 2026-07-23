package run

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

type mergeTestEvent struct{}

func (mergeTestEvent) EventName() string { return "MergeTest" }

type mergeTestHooks struct {
	name   string
	detect hookkit.DetectFunc
	codec  *hookkit.Codec
	apply  func(*hookkit.Dialect)
}

func (c mergeTestHooks) Contribute(reg Registry) {
	d := reg.Ensure(c.name, c.detect, c.codec)
	c.apply(d)
}

func TestContribute_MergesSameDialectHooks(t *testing.T) {
	var first, second atomic.Int32
	codec := hookkit.NewCodec("merge", fmt.Errorf("empty"), fmt.Errorf("decode"), fmt.Errorf("name required"))
	codec.Register("MergeTest", func([]byte) (hookkit.Event, error) {
		return mergeTestEvent{}, nil
	})
	detect := func([]byte) bool { return true }

	h1 := mergeTestHooks{
		name: "merge", detect: detect, codec: codec,
		apply: func(d *hookkit.Dialect) {
			d.Register(hookkit.ObserveHandler(func(context.Context, mergeTestEvent) error {
				first.Add(1)
				return nil
			}))
		},
	}
	h2 := mergeTestHooks{
		name: "merge", detect: detect, codec: codec,
		apply: func(d *hookkit.Dialect) {
			d.Register(hookkit.ObserveHandler(func(context.Context, mergeTestEvent) error {
				second.Add(1)
				return nil
			}))
		},
	}

	router := newRouter()
	h1.Contribute(router)
	h2.Contribute(router)

	code := serve(context.Background(), router, strings.NewReader(`{"hook_event_name":"MergeTest"}`), &bytes.Buffer{}, &bytes.Buffer{})
	if code != 0 {
		t.Fatalf("exit = %d", code)
	}
	if first.Load() != 1 || second.Load() != 1 {
		t.Fatalf("handler calls first=%d second=%d, want 1 each", first.Load(), second.Load())
	}
}
