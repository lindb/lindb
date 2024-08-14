package matching

type Property[S any, C any, R any] struct {
	Fn   func(context C, source S) R
	Name string
}

func PropertyOf[S any, C any, R any](name string, fn func(context C, source S) R) *Property[S, C, R] {
	return &Property[S, C, R]{
		Name: name,
		Fn:   fn,
	}
}

func (p *Property[S, C, R]) Matching(pattern Pattern) *PropertyPattern[S, C, R] {
	return &PropertyPattern[S, C, R]{
		Property: p,
		Pattern:  pattern,
	}
}
