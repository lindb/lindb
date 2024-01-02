package tree

import "github.com/lindb/lindb/pkg/collections"

type StackableVisitorContext[C any] struct {
	stack   *collections.Stack
	context C
}

func NewStackableVisitorContext[C any](context C) *StackableVisitorContext[C] {
	return &StackableVisitorContext[C]{
		stack:   collections.NewStack(),
		context: context,
	}
}

func (c *StackableVisitorContext[C]) GetContext() (r C) {
	return c.context
}

func (c *StackableVisitorContext[C]) Push(node Node) {
	c.stack.Push(node)
}

func (c *StackableVisitorContext[C]) Pop() {
	_ = c.stack.Pop()
}

func (c *StackableVisitorContext[C]) GetPreviousNode() Node {
	if c.stack.Size() > 1 {
		return c.stack.Get(1).(Node)
	}
	return nil
}

type StackableAstVisitor[C any] struct{}

func (v *StackableAstVisitor[C]) Visit(context any, node Node) any {
	stackCxt := context.(*StackableVisitorContext[C])
	stackCxt.Push(node)
	defer func() {
		stackCxt.Pop()
	}()
	// FIXME: ??
	return nil
}
