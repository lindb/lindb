package printer

import (
	"regexp"
	"strings"
)

const (
	VerticalLine     = "│"
	LastNode         = "└─"
	IntermediateNode = "├─"
)

type Render interface {
	Render(paln *PlanRepresentation) string
}

type TextRender struct {
	level int
}

func NewTextRender(level int) Render {
	return &TextRender{
		level: level,
	}
}

func (r *TextRender) Render(plan *PlanRepresentation) string {
	root := plan.getRoot()
	sb := &strings.Builder{}
	hasChildren := hasChildren(plan, root)

	return strings.TrimSuffix(r.writeTextOutput(sb, plan, NewIndent(r.level, hasChildren), root), "\n")
}

func (r *TextRender) writeTextOutput(output *strings.Builder, plan *PlanRepresentation, indent *Indent, node *NodeRepresentation) string {
	output.WriteString(indent.nodeIndent())
	output.WriteString(node.getName())
	var kvs []string
	for key, value := range node.descriptor {
		kvs = append(kvs, key+" = "+value)
	}
	output.WriteString("[" + strings.Join(kvs, ", ") + "]")
	output.WriteString("\n")

	output.WriteString(indentMultilineString("Layout: "+formatSymbols(node.outputs), indent.detailIndent()))
	output.WriteString("\n")

	if len(node.details) > 0 {
		details := strings.Join(node.details, "\n")
		details = indentMultilineString(details, indent.detailIndent())
		output.WriteString(details)
		if !strings.HasSuffix(details, "\n") {
			output.WriteString("\n")
		}
	}

	// process children
	childrenIDs := node.children
	for i, childID := range childrenIDs {
		child := plan.getNode(childID)
		if child != nil {
			r.writeTextOutput(output, plan,
				indent.forChild(i == len(childrenIDs)-1, hasChildren(plan, child)), child)
		}
	}

	return output.String()
}

func indentMultilineString(str, indent string) string {
	m1 := regexp.MustCompile("(?m)^")
	return m1.ReplaceAllString(str, indent)
}

type Indent struct {
	firstLinePrefix string
	nextLinesPrefix string
	hasChildren     bool
}

func NewIndent(level int, hasChildren bool) *Indent {
	indent := indentString(level)
	return &Indent{
		firstLinePrefix: indent,
		nextLinesPrefix: indent,
		hasChildren:     hasChildren,
	}
}

func (i *Indent) nodeIndent() string {
	return i.firstLinePrefix
}

func (i *Indent) detailIndent() string {
	indent := ""
	if i.hasChildren {
		indent = VerticalLine
	}
	return i.nextLinesPrefix + pad(indent, 2)
}

func (i *Indent) forChild(last, hasChildren bool) *Indent {
	var (
		first string
		next  string
	)
	if last {
		first = pad(LastNode, 3)
		next = pad("", 3)
	} else {
		first = pad(IntermediateNode, 3)
		next = pad(VerticalLine, 3)
	}
	return &Indent{
		firstLinePrefix: i.nextLinesPrefix + first,
		nextLinesPrefix: i.nextLinesPrefix + next,
		hasChildren:     hasChildren,
	}
}

func indentString(indent int) string {
	return strings.Repeat("    ", indent)
}

func pad(text string, length int) string {
	return text + strings.Repeat(" ", length-len([]rune(text)))
}

func hasChildren(plan *PlanRepresentation, node *NodeRepresentation) bool {
	for _, childID := range node.children {
		child := plan.getNode(childID)
		if child != nil {
			return true
		}
	}
	return false
}
