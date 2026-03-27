package template

// TemplateBindings maps placeholder keys to values for argv templating.
// TemplateValue reports whether key is defined; value may be empty when ok is true.
type TemplateBindings interface {
	TemplateValue(key string) (value string, ok bool)
}
