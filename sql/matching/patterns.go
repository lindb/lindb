package matching

import (
	"fmt"
	"reflect"

	"github.com/samber/lo"
)

func TypeOf(expectedType reflect.Type) *Pattern {
	return &Pattern{
		Accept: func(context, val any, captures *Captures) []*Match {
			fmt.Printf("check type %v=%v\n", reflect.TypeOf(val), expectedType)
			if reflect.TypeOf(val) == expectedType {
				return []*Match{NewMatch(captures)}
			}
			return nil
		},
	}
}

func CapturedAs(capture *Capture, previous *Pattern) *Pattern {
	fmt.Printf("add captured..................%v,%p\n", capture, capture)
	return &Pattern{
		Previous: previous,
		Accept: func(context, val any, captures *Captures) []*Match {
			newCaptures := captures.AddAll(OfNullable(capture, val))
			return []*Match{NewMatch(newCaptures)}
		},
	}
}

type Pattern struct {
	Accept   func(context any, val any, captures *Captures) []*Match
	Previous *Pattern
}

//	func With(pattern *PropPattern, previous Pattern) Pattern {
//		return &WithPattern{
//			previous:    previous,
//			propPattern: pattern,
//		}
//	}

func (p *Pattern) Match(context any, val any, captures *Captures) []*Match {
	if p.Previous != nil {
		matches := lo.FlatMap(p.Previous.Match(context, val, captures), func(match *Match, index int) []*Match {
			a := p.Accept(context, val, match.Captures)
			fmt.Printf("flat map aaaa=%v,%T,%v\n", a, val, val)
			return a
		})
		fmt.Printf("pattern match.......%v\n", matches)
		return matches
	}
	return p.Accept(context, val, captures)
}

// type TypeOfPattern struct {
// 	expectedType reflect.Type
// }
//
// func (p *TypeOfPattern) Accept(context any, val any, captures *Captures) []*Match {
// 	if reflect.TypeOf(val) == p.expectedType {
// 		return []*Match{NewMatch(captures)}
// 	}
// 	return nil
// }
//
// func (p *TypeOfPattern) Previous() Pattern {
// 	return nil
// }

type CapturePattern struct {
	previous *Pattern
	capture  *Capture
}

func (p *CapturePattern) Accept(context any, val any, captures *Captures) []*Match {
	newCaptures := captures.AddAll(OfNullable(p.capture, val))
	return []*Match{NewMatch(newCaptures)}
}

type WithPattern struct {
	previous Pattern
	// propPattern *PropertyPattern
}

func (p *WithPattern) Accept(context any, val any, captures *Captures) []*Match {
	return nil
}

func (p *WithPattern) Previous() Pattern {
	return p.previous
}

type PropertyPattern[S any, C any, R any] struct {
	Property *Property[S, C, R]
	Pattern  Pattern
}
