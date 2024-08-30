package tree

import "fmt"

type DefaultTraversalVisitor struct {
	Process func(n Node)
}

func (v *DefaultTraversalVisitor) Visit(context any, n Node) (r any) {
	fmt.Printf("express visit = %T value=%v\n", n, n)
	switch node := n.(type) {
	case *ArithmeticBinaryExpression:
		node.Left.Accept(context, v)
		node.Right.Accept(context, v)
	default:
		// TODO: remove
		fmt.Printf("default traversal visitor not support..................=%T\n", n)
	}
	if v.Process != nil {
		// if visitor has process func, invoke it
		v.Process(n)
	}
	return
}
