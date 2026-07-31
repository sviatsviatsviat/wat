package run

import (
	"fmt"
	"io"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// serveHints holds optional install-time dispatch overrides from process args.
type serveHints struct {
	agent string
	event string
}

func parseServeHints(args []string) serveHints {
	var h serveHints
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--agent":
			if i+1 < len(args) {
				i++
				h.agent = args[i]
			}
		case "--event":
			if i+1 < len(args) {
				i++
				h.event = args[i]
			}
		}
	}
	return h
}

func resolveDialect(router *router, raw []byte, hints serveHints, errw io.Writer) (name string, d *hookkit.Dialect, code int) {
	if hints.agent != "" {
		var ok bool
		d, ok = router.lookup(hints.agent)
		if !ok || d == nil || d.Codec() == nil {
			_, _ = fmt.Fprintf(errw, "run: unknown dialect %q\n", hints.agent)
			return "", nil, 1
		}
		name = hints.agent
		if detected, _, ok := router.detect(raw); ok && detected != "" && detected != name {
			_, _ = fmt.Fprintf(errw, "run: warning: --agent %q disagrees with detected dialect %q; using --agent\n", name, detected)
		}
		return name, d, 0
	}
	var ok bool
	name, d, ok = router.detect(raw)
	if !ok || d == nil || d.Codec() == nil {
		_, _ = fmt.Fprintln(errw, "run: unknown dialect")
		return "", nil, 1
	}
	return name, d, 0
}

func resolveEvent(d *hookkit.Dialect, dialectName string, raw []byte, hints serveHints, errw io.Writer) (eventName string, code int) {
	if hints.event != "" {
		eventName = hints.event
		if peeked, err := hookkit.PeekHookEventName(raw); err == nil && peeked != "" && peeked != eventName {
			_, _ = fmt.Fprintf(errw, "run: warning: --event %q disagrees with hook_event_name %q; using --event\n", eventName, peeked)
		}
		return eventName, 0
	}
	eventName, err := d.Codec().EventName(raw)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", dialectName, err)
		return "", 1
	}
	return eventName, 0
}
