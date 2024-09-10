package matching

import (
	"fmt"

	"go.uber.org/atomic"
)

// TODO: remove this package
var (
	EmptyCaptures = &Captures{}
	sequence      = atomic.NewUint64(0)
)

type Capture struct {
	seq uint64
}

func NewCapture() *Capture {
	return &Capture{
		seq: sequence.Inc(),
	}
}

type Captures struct {
	Capture *Capture
	Val     any
	Tail    *Captures
}

func OfNullable(capture *Capture, val any) *Captures {
	if capture == nil {
		return EmptyCaptures
	}
	return &Captures{
		Capture: capture,
		Val:     val,
		Tail:    EmptyCaptures,
	}
}

func (c *Captures) Get(capture *Capture) any {
	if c == EmptyCaptures {
		panic("unknow capture, was it registerd in the pattern")
	}
	if c.Capture == capture {
		fmt.Printf("get capture vall....,%T,tail=%T,a=%p,b=%p\n", c.Val, c.Tail.Val, c.Capture, capture)
		return c.Val
	}
	fmt.Printf("get capture vall.... from tail%T\n", c.Val)
	return c.Tail.Get(capture)
}

func (c *Captures) AddAll(other *Captures) *Captures {
	if c == EmptyCaptures {
		return other
	}
	return &Captures{
		Capture: c.Capture,
		Val:     c.Val,
		Tail:    c.Tail.AddAll(other),
	}
}
