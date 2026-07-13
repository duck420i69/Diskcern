package providers

type Registry struct {
	providers []Provider
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) Register(p Provider) {
	r.providers = append(r.providers, p)
}

func (r *Registry) Detect(path string, isDir bool) (Provider, ScanDirective) {
	for _, p := range r.providers {
		if matched, directive := p.Detect(path, isDir); matched {
			return p, directive
		}
	}
	return nil, ContinueTraversal
}

// GlobalProviders returns providers that don't rely on path detection (like Docker)
func (r *Registry) GlobalProviders() []Provider {
	var globals []Provider
	for _, p := range r.providers {
		if matched, _ := p.Detect("", false); matched {
			globals = append(globals, p)
		}
	}
	return globals
}
