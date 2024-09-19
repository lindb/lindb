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
		return v.formatBinaryExpression(string(node.Operator), node.Left, node.Right)
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
	case *LongLiteral:
		return fmt.Sprintf("%d", node.Value)
	case *Constant:
		// TODO: add more
		return fmt.Sprintf("%v", node.Value)
	// case *ArithmeticBinaryExpression:
	// 	return v.formatBinaryExpression(string(node.Operator), node.Left, node.Right)
	case *Call:
		var args []string
		for _, arg := range node.Args {
			args = append(args, arg.Accept(context, v).(string))
		}
		return fmt.Sprintf("%v(%s)", node.Function, strings.Join(args, ","))
	case *FunctionCall:
		var args []string
		for _, arg := range node.Arguments {
			args = append(args, arg.Accept(context, v).(string))
		}
		return fmt.Sprintf("%v(%s)", node.Name, strings.Join(args, ","))
	case *Cast:
		return fmt.Sprintf("CAST(%v as %s)", node.Expression.Accept(context, v), node.Type)
	default:
		panic(fmt.Sprintf("expression formatter unsupport node:%T", n))
	}
}

func (v *FormatVisitor) visitDereferenceExpression(context any, node *DereferenceExpression) (r any) {
	field := "*" // FIXME: all select?
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

func (v *FormatVisitor) formatBinaryExpression(operator string, left, right Expression) any {
	return fmt.Sprintf("(%v %s %v)", left.Accept(nil, v), operator, right.Accept(nil, v))
}

func (v *FormatVisitor) formatStringLiteral(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "''"))
}

func (v *FormatVisitor) formatIdentifier(s string) string {
	return strings.ReplaceAll(s, "\"", "\"\"")
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
