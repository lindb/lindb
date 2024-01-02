package tree

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/lindb/lindb/sql/grammar"
)

var keyWords = make(map[string]string)

func init() {
	typ := reflect.TypeOf(&grammar.NonReservedContext{})
	for i := 0; i < typ.NumMethod(); i++ {
		methodName := typ.Method(i).Name
		if isUpper(methodName) {
			keyWords[methodName] = methodName
		}
	}
}

func FormatExpression(expression Expression) string {
	msg := expression.Accept(nil, NewFormatVisitor())
	if msg != nil {
		return msg.(string)
	}
	return ""
}

type FormatVisitor struct{}

func NewFormatVisitor() Visitor {
	return &FormatVisitor{}
}

func (v *FormatVisitor) Visit(context any, n Node) any {
	switch node := n.(type) {
	case *ComparisonExpression:
		return v.formatBinaryExpression(node.Operator, node.Left, node.Right)
	case *DereferenceExpression:
		return v.visitDereferenceExpression(context, node)
	case *SymbolReference:
		return v.formatIdentifier(node.Name)
	case *FieldReference:
		// add colon so this won't parse
		return fmt.Sprintf(":input(%d)", node.FieldIndex)
	case *Identifier:
		return v.visitIdentifier(context, node)
	case *StringLiteral:
		return v.formatStringLiteral(node.Value)
	case *Cast:
		return "Cast need impl"
	default:
		panic(fmt.Sprintf("expression formatter unsupport node:%T", n))
	}
}

func (v *FormatVisitor) visitDereferenceExpression(context any, node *DereferenceExpression) (r any) {
	field := "*"
	if node.Field != nil {
		fieldMsg := node.Field.Accept(context, v)
		if fieldMsg != nil {
			field = fieldMsg.(string)
		}
	}
	return fmt.Sprintf("%v.%s", node.Base.Accept(context, v), field)
}

func (v *FormatVisitor) visitIdentifier(context any, node *Identifier) (r any) {
	if node.Delimited || reserved(node.Value) {
		return v.formatIdentifier(node.Value)
	}
	return node.Value
}

func (v *FormatVisitor) formatBinaryExpression(operator ComparisonOperator, left, right Expression) any {
	return fmt.Sprintf("(%v %s %v)", left.Accept(v, nil), operator, right.Accept(v, nil))
}

func (v *FormatVisitor) formatStringLiteral(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "''"))
}

func (v *FormatVisitor) formatIdentifier(s string) string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(s, "\"", "\"\""))
}

func reserved(name string) bool {
	_, ok := keyWords[strings.ToUpper(name)]
	return ok
}

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) {
			return false
		}
	}
	return true
}
