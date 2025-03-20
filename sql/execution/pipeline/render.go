package pipeline

import (
	"fmt"
	"strings"

	"github.com/lindb/lindb/pkg/strutil"
	"github.com/lindb/lindb/sql/execution/pipeline/operator"
)

func renderText(root operator.Operator) string {
	sb := &strings.Builder{}
	return strings.TrimSuffix(writeTextOutput(sb, strutil.NewIndent(0, len(root.Children()) > 0), root), "\n")
}

func writeTextOutput(sb *strings.Builder, indent *strutil.Indent, node operator.Operator) string {
	sb.WriteString(indent.NodeIndent())
	sb.WriteString(strings.TrimPrefix(fmt.Sprintf("%T\n", node), "*"))
	children := node.Children()
	for i, child := range children {
		writeTextOutput(sb, indent.ForChild(i == len(children)-1, len(child.Children()) > 0), child)
	}
	return sb.String()
}
