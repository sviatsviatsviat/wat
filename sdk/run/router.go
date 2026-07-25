package run

import (
	"sort"

	"github.com/sviatsviatsviat/wat/internal/hookkit"
)

type routerEntry struct {
	detect  hookkit.DetectFunc
	dialect *hookkit.Dialect
}

// router selects an agent dialect for one hook serve cycle.
// Ensure must finish before detect; not safe for concurrent use.
type router struct {
	order  []string
	byName map[string]routerEntry
}

func newRouter() *router {
	return &router{
		byName: make(map[string]routerEntry),
	}
}

// Ensure registers detect and codec for name when missing and returns the dialect.
// Same name returns the existing dialect so handlers can append.
func (r *router) Ensure(name string, detect hookkit.DetectFunc, codec *hookkit.Codec) *hookkit.Dialect {
	if name == "" {
		panic("run: Ensure: empty name")
	}
	if codec == nil {
		panic("run: Ensure: nil codec")
	}
	if e, ok := r.byName[name]; ok {
		return e.dialect
	}
	d := hookkit.NewDialect(codec)
	r.order = append(r.order, name)
	r.byName[name] = routerEntry{detect: detect, dialect: d}
	return d
}

func (r *router) registry() {}

func (r *router) detect(raw []byte) (name string, d *hookkit.Dialect, ok bool) {
	for _, name := range r.order {
		e := r.byName[name]
		if e.detect != nil && e.detect(raw) {
			return name, e.dialect, true
		}
	}
	return "", nil, false
}

func (r *router) manifest() Manifest {
	manifest := Manifest{Version: ManifestVersion}
	for name, entry := range r.byName {
		for _, event := range entry.dialect.EventNames() {
			manifest.Registrations = append(manifest.Registrations, Registration{
				Dialect:      name,
				Event:        event,
				HandlerCount: len(entry.dialect.HandlersFor(event)),
			})
		}
	}
	sort.Slice(manifest.Registrations, func(i, j int) bool {
		left, right := manifest.Registrations[i], manifest.Registrations[j]
		if left.Dialect != right.Dialect {
			return left.Dialect < right.Dialect
		}
		return left.Event < right.Event
	})
	return manifest
}
