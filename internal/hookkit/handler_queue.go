package hookkit

// HandlerQueue holds dialect handler registrations until Install applies them.
// Agent SDKs embed or wrap HandlerQueue behind a typed fluent UseHooks API.
type HandlerQueue struct {
	pending []func(*Dialect)
}

// DialectEnsurer ensures a named dialect and returns its handler bag.
// run.Registry satisfies DialectEnsurer.
type DialectEnsurer interface {
	Ensure(name string, detect DetectFunc, codec *Codec) *Dialect
}

// Bind queues register(d, fn) for the next Install. A nil queue is a no-op.
func Bind[F any](q *HandlerQueue, fn F, register func(*Dialect, F)) *HandlerQueue {
	if q == nil {
		return nil
	}
	q.pending = append(q.pending, func(d *Dialect) { register(d, fn) })
	return q
}

// Install ensures the dialect on ensurer and runs queued registrations in order.
// A nil receiver or ensurer is a no-op.
func (q *HandlerQueue) Install(ensurer DialectEnsurer, name string, detect DetectFunc, codec *Codec) {
	if q == nil || ensurer == nil {
		return
	}
	d := ensurer.Ensure(name, detect, codec)
	for _, apply := range q.pending {
		apply(d)
	}
}
