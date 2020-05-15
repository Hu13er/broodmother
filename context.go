package broodmother

type Context interface {
	Global() Context
	Parent() Context

	Imports() map[string]string
	Package() string
	Path() string

	Get(key interface{}) (interface{}, bool)
	Set(key, value interface{})
}

type background struct {
	imports map[string]string
	pkg     string
	path    string
	values  map[interface{}]interface{}
}

func newBackground(pkg, path string, imports map[string]string) Context {
	return &background{
		imports: imports,
		pkg:     pkg,
		path:    path,
		values:  make(map[interface{}]interface{}),
	}
}

func (b *background) Package() string            { return b.pkg }
func (b *background) Imports() map[string]string { return b.imports }
func (b *background) Global() Context            { return b }
func (b *background) Parent() Context            { return nil }
func (b *background) Path() string               { return b.path }

func (b *background) Get(key interface{}) (interface{}, bool) {
	if b.values == nil {
		return nil, false
	}
	value, exists := b.values[key]
	return value, exists
}

func (b *background) Set(key, value interface{}) {
	b.values[key] = value
}

type layer struct {
	Context
	values map[interface{}]interface{}
}

func newLayer(parent Context) Context {
	return &layer{
		Context: parent,
		values:  make(map[interface{}]interface{}),
	}
}

func (l *layer) Parent() Context {
	return l.Context
}

func (l *layer) Set(key, value interface{}) {
	l.values[key] = value
}

func (l *layer) Get(key interface{}) (interface{}, bool) {
	if l.values == nil {
		return l.Context.Get(key)
	}
	value, exists := l.values[key]
	if !exists {
		return l.Context.Get(key)
	}
	return value, true
}
