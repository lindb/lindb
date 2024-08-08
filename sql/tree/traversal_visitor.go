package tree

import "fmt"

type DefaultTraversalVisitor struct{}

func (v *DefaultTraversalVisitor) Visit(context any, node Node) (r any) {
	fmt.Printf("express visis = %T value=%v\n", node, node)
	return
}
