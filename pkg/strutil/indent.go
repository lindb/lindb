package strutil

import (
	"strings"

	"github.com/lindb/lindb/constants"
)

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

func (i *Indent) NodeIndent() string {
	return i.firstLinePrefix
}

func (i *Indent) DetailIndent() string {
	indent := ""
	if i.hasChildren {
		indent = constants.VerticalLine
	}
	return i.nextLinesPrefix + pad(indent, 2)
}

func (i *Indent) ForChild(last, hasChildren bool) *Indent {
	var (
		first string
		next  string
	)
	if last {
		first = pad(constants.LastNode, 3)
		next = pad("", 3)
	} else {
		first = pad(constants.IntermediateNode, 3)
		next = pad(constants.VerticalLine, 3)
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
