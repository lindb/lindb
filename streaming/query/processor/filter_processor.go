package processor

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/samber/lo"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/sql/tree"
	"github.com/lindb/lindb/streaming/query/context"
)

type FilterProcessor struct {
	ctx *context.QueryContext

	node *plan.FilterNode

	child Processor

	inbound *Queue
}

func NewFilterProcessor(ctx *context.QueryContext,
	node *plan.FilterNode,
	child Processor,
) Processor {
	return &FilterProcessor{
		ctx:     ctx,
		node:    node,
		child:   child,
		inbound: NewQueue(make(chan any)),
	}
}

func (f *FilterProcessor) Run(output chan<- any) {
	v := &visitor{}
	for {
		source, ok := f.inbound.Consume(f.ctx.Context)
		if !ok {
			break
		}
		if val, ok := v.Visit(source, f.node.Predicate).(bool); val && ok {
			output <- source
		}
	}
}

func (f *FilterProcessor) GetInbounds() []chan any {
	return []chan any{f.inbound.GetInbound()}
}

func (f *FilterProcessor) Children() []Processor {
	return []Processor{f.child}
}

func (f *FilterProcessor) String() string {
	return fmt.Sprintf("Filter[predicate=%s]", tree.FormatExpression(f.node.Predicate))
}

type visitor struct{}

func (v *visitor) Visit(event any, n tree.Node) any {
	switch node := n.(type) {
	case *tree.ComparisonExpression:
		v, err := GetFieldValue(event, getValue(node.Left))
		if err != nil {
			fmt.Println(err)
			return false
		}
		return v == getValue(node.Right)
	case *tree.InPredicate:
		var values []string
		if inListExpression, ok := node.ValueList.(*tree.InListExpression); ok {
			values = lo.Map(inListExpression.Values, func(item tree.Expression, index int) string {
				return getValue(item)
			})
		}
		v, err := GetFieldValue(event, getValue(node.Value))
		if err != nil {
			fmt.Println(err)
			return false
		}
		return lo.Contains(values, v.(string))
	case *tree.LogicalExpression:
		for _, term := range node.Terms {
			val, ok := term.Accept(event, v).(bool)
			if !ok {
				return false
			}
			if node.Operator == tree.LogicalOR && val {
				return true
			} else if !val {
				return false
			}
		}
		return true
	default:
		panic(fmt.Errorf("not support,%T", n))
	}
}

func GetFieldValue(s interface{}, field string) (interface{}, error) {
	v := reflect.ValueOf(s)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("input is not a struct or pointer to struct")
	}

	f := v.FieldByName(lo.PascalCase(field))
	if !f.IsValid() {
		return nil, fmt.Errorf("no such field: %s in struct", field)
	}

	if !f.CanInterface() {
		return nil, fmt.Errorf("cannot access field %s", field)
	}

	return f.Interface(), nil
}

func getValue(n tree.Expression) string {
	switch node := n.(type) {
	case *tree.StringLiteral:
		return node.Value
	case *tree.Identifier:
		return node.Value
	}
	return ""
}
