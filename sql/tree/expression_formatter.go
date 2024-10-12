package tree

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/samber/lo"

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
	case *NotExpression:
		return v.formatNotExpression(node)
	case *InPredicate:
		return v.formatInPredicate(node)
	case *InListExpression:
		return v.formatInListExpression(node)
	case *LikePredicate:
		return v.formatLikePredicate(node)
	case *RegexPredicate:
		return v.formatRegexPredicate(node)
	case *TimePredicate:
		return v.formatTimestampPredicate(node)
	case *LogicalExpression:
		return v.formatLogical(node)
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
	case *IntervalLiteral:
		return v.formatIntervalLiteral(node)
	case *Constant:
		// TODO: add more
		return fmt.Sprintf("%v", node.Value)
	// case *ArithmeticBinaryExpression:
	// 	return v.formatBinaryExpression(string(node.Operator), node.Left, node.Right)
	case *FunctionCall:
		return fmt.Sprintf("%v(%s)", node.Name, strings.Join(lo.Map(node.Arguments, func(arg Expression, index int) string {
			return arg.Accept(context, v).(string)
		}), ","))
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

func (v *FormatVisitor) formatBinaryExpression(operator string, left, right Expression) string {
	return fmt.Sprintf("(%v %s %v)", left.Accept(nil, v), operator, right.Accept(nil, v))
}

func (v *FormatVisitor) formatNotExpression(node *NotExpression) string {
	return fmt.Sprintf("(NOT %s)", node.Value.Accept(nil, v))
}

func (v *FormatVisitor) formatInPredicate(node *InPredicate) string {
	return fmt.Sprintf("(%v IN %v)", node.Value.Accept(nil, v), node.ValueList.Accept(nil, v))
}

func (v *FormatVisitor) formatInListExpression(node *InListExpression) string {
	return fmt.Sprintf("(%v)", strings.Join(lo.Map(node.Values, func(item Expression, index int) string {
		return item.Accept(nil, v).(string)
	}), ","))
}

func (v *FormatVisitor) formatLikePredicate(node *LikePredicate) string {
	return fmt.Sprintf("(%s LIKE %s)", node.Value.Accept(nil, v), node.Pattern.Accept(nil, v))
}

func (v *FormatVisitor) formatRegexPredicate(node *RegexPredicate) string {
	return fmt.Sprintf("(%s REGEXP %s)", node.Value.Accept(nil, v), node.Pattern.Accept(nil, v))
}

func (v *FormatVisitor) formatTimestampPredicate(node *TimePredicate) string {
	return fmt.Sprintf("(time %s %v)", node.Operator, node.Value.Accept(nil, v))
}

func (v *FormatVisitor) formatStringLiteral(s string) string {
	return fmt.Sprintf("'%s'", strings.ReplaceAll(s, "'", "''"))
}

func (v *FormatVisitor) formatIntervalLiteral(node *IntervalLiteral) string {
	return fmt.Sprintf("INTERVAL %v %s", node.Value, node.Unit)
}

func (v *FormatVisitor) formatIdentifier(s string) string {
	return strings.ReplaceAll(s, "\"", "\"\"")
}

func (v *FormatVisitor) formatLogical(node *LogicalExpression) string {
	return fmt.Sprintf("(%v)", strings.Join(lo.Map(node.Terms, func(item Expression, index int) string {
		return v.Visit(nil, item).(string)
	}), fmt.Sprintf(" %v ", node.Operator)))
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
