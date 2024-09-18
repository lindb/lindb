package tree

import "fmt"

type DefaultTraversalVisitor struct {
	PreProcess  func(n Node)
	PostProcess func(n Node)
}

func (v *DefaultTraversalVisitor) Visit(context any, n Node) (r any) {
	fmt.Printf("express visit = %T value=%v\n", n, n)
	if v.PreProcess != nil {
		// do pre process if has pre func
		v.PreProcess(n)
	}
	switch node := n.(type) {
	case *ArithmeticBinaryExpression:
		_ = node.Left.Accept(context, v)
		_ = node.Right.Accept(context, v)
	case *FunctionCall:
		for _, arg := range node.Arguments {
			_ = arg.Accept(context, v)
		}
	default:
		// TODO: remove
		fmt.Printf("default traversal visitor not support..................=%T\n", n)
	}
	if v.PostProcess != nil {
		// do post process if has post func
		v.PostProcess(n)
	}
	return
}
