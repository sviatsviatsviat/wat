package hookkit

// DetectFunc reports whether raw matches a dialect.
type DetectFunc func(raw []byte, getenv func(string) string) bool

type routerEntry struct {
	detect  DetectFunc
	dialect *Dialect
}

// Router selects an agent dialect for one hook serve cycle.
// Ensure must finish before Detect; not safe for concurrent use.
type Router struct {
	order  []string
	byName map[string]routerEntry
}

var defaultRouter = NewRouter()

// NewRouter returns an empty dialect router.
func NewRouter() *Router {
	return &Router{
		byName: make(map[string]routerEntry),
	}
}

// DefaultRouter returns the process-wide router used by run.Main.
func DefaultRouter() *Router {
	return defaultRouter
}

// Ensure registers d under name when missing and returns the dialect for name.
func (r *Router) Ensure(name string, detect DetectFunc, d *Dialect) *Dialect {
	if name == "" {
		panic("hookkit: Ensure: empty name")
	}
	if d == nil {
		panic("hookkit: Ensure: nil dialect")
	}
	if e, ok := r.byName[name]; ok {
		return e.dialect
	}
	r.order = append(r.order, name)
	r.byName[name] = routerEntry{detect: detect, dialect: d}
	return d
}

// Detect selects a dialect for raw, or forced when non-empty.
func (r *Router) Detect(raw []byte, getenv func(string) string, forced string) (name string, d *Dialect, ok bool) {
	if forced != "" {
		e, ok := r.byName[forced]
		if !ok {
			return "", nil, false
		}
		return forced, e.dialect, true
	}
	for _, name := range r.order {
		e := r.byName[name]
		if e.detect != nil && e.detect(raw, getenv) {
			return name, e.dialect, true
		}
	}
	return "", nil, false
}
