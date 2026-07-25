package hookkit

import (
	"sync/atomic"
	"testing"
)

type fakeEnsurer struct {
	name   string
	detect DetectFunc
	codec  *Codec
	d      *Dialect
}

func (f *fakeEnsurer) Ensure(name string, detect DetectFunc, codec *Codec) *Dialect {
	if f.d == nil {
		f.name = name
		f.detect = detect
		f.codec = codec
		f.d = &Dialect{}
	}
	return f.d
}

func TestHandlerQueue_InstallAppliesInOrder(t *testing.T) {
	var order []int
	var q HandlerQueue
	Bind(&q, 1, func(_ *Dialect, n int) { order = append(order, n) })
	Bind(&q, 2, func(_ *Dialect, n int) { order = append(order, n) })

	ensurer := &fakeEnsurer{}
	q.Install(ensurer, "test", func([]byte) bool { return true }, nil)

	if ensurer.name != "test" {
		t.Fatalf("Ensure name = %q, want test", ensurer.name)
	}
	if got, want := order, []int{1, 2}; len(got) != 2 || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("order = %v, want %v", got, want)
	}
}

func TestHandlerQueue_NilReceiverAndEnsurerNoop(t *testing.T) {
	var called atomic.Bool
	var q *HandlerQueue
	if got := Bind(q, 1, func(*Dialect, int) { called.Store(true) }); got != nil {
		t.Fatalf("Bind(nil) = %v, want nil", got)
	}
	q.Install(&fakeEnsurer{}, "x", nil, nil)

	q2 := &HandlerQueue{}
	Bind(q2, 1, func(*Dialect, int) { called.Store(true) })
	q2.Install(nil, "x", nil, nil)
	if called.Load() {
		t.Fatal("Install(nil ensurer) applied registrations")
	}
}
