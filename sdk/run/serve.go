package run

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

// Serve merges hooks into a local dialect router (same dialect name appends
// handlers), runs one hook dispatch cycle on os.Stdin / os.Stdout / os.Stderr,
// then os.Exit with the resulting code.
//
// Optional --agent and --event flags in os.Args force dialect and event
// selection (skipping payload detect/peek). When a hint disagrees with the
// payload, Serve warns on stderr and continues with the hint.
func Serve(hooks ...Hooks) {
	r := newRouter()
	contribute(r, hooks)
	code := serve(context.Background(), r, os.Stdin, os.Stdout, os.Stderr, parseServeHints(os.Args[1:]))
	os.Exit(code)
}

func contribute(r Registry, hooks []Hooks) {
	for _, h := range hooks {
		if h == nil {
			continue
		}
		h.Contribute(r)
	}
}

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

func serve(ctx context.Context, router *router, in io.Reader, out io.Writer, errw io.Writer, hints serveHints) int {
	raw, err := io.ReadAll(in)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: read stdin: %v\n", err)
		return 1
	}
	if len(raw) == 0 {
		_, _ = fmt.Fprintln(errw, "run: empty stdin")
		return 1
	}

	var name string
	var d *hookkit.Dialect
	if hints.agent != "" {
		var ok bool
		d, ok = router.lookup(hints.agent)
		if !ok || d == nil || d.Codec() == nil {
			_, _ = fmt.Fprintf(errw, "run: unknown dialect %q\n", hints.agent)
			return 1
		}
		name = hints.agent
		if detected, _, ok := router.detect(raw); ok && detected != "" && detected != name {
			_, _ = fmt.Fprintf(errw, "run: warning: --agent %q disagrees with detected dialect %q; using --agent\n", name, detected)
		}
	} else {
		var ok bool
		name, d, ok = router.detect(raw)
		if !ok || d == nil || d.Codec() == nil {
			_, _ = fmt.Fprintln(errw, "run: unknown dialect")
			return 1
		}
	}

	var eventName string
	if hints.event != "" {
		eventName = hints.event
		if peeked, err := hookkit.PeekHookEventName(raw); err == nil && peeked != "" && peeked != eventName {
			_, _ = fmt.Fprintf(errw, "run: warning: --event %q disagrees with hook_event_name %q; using --event\n", eventName, peeked)
		}
	} else {
		eventName, err = d.Codec().EventName(raw)
		if err != nil {
			_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", name, err)
			return 1
		}
	}

	handlers := d.HandlersFor(eventName)
	if len(handlers) == 0 {
		return 0
	}

	event, err := d.Codec().DecodeAs(raw, eventName)
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: decode: %v\n", name, err)
		return 1
	}

	var acc hookkit.Output
	for _, h := range handlers {
		outVal, err := h.Invoke(ctx, event)
		if err != nil {
			_, _ = fmt.Fprintf(errw, "run: %s: handler: %v\n", name, err)
			return 1
		}
		if outVal == nil || outVal.IsZero() {
			continue
		}
		if acc == nil || acc.IsZero() {
			acc = outVal
		} else {
			merged, warnings, err := acc.Merge(outVal)
			if err != nil {
				_, _ = fmt.Fprintf(errw, "run: %s: merge: %v\n", name, err)
				return 1
			}
			for _, w := range warnings {
				_, _ = fmt.Fprintf(errw, "run: %s: merge: %s\n", name, w)
			}
			acc = merged
		}
		if acc.Stop() {
			break
		}
	}

	if acc == nil || acc.IsZero() {
		return 0
	}

	stdout, exitCode, err := acc.Encode()
	if err != nil {
		_, _ = fmt.Fprintf(errw, "run: %s: encode: %v\n", name, err)
		return 1
	}
	if len(stdout) > 0 {
		if _, err := out.Write(stdout); err != nil {
			_, _ = fmt.Fprintf(errw, "run: write stdout: %v\n", err)
			return 1
		}
	}
	return exitCode
}
