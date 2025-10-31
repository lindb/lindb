package streaming

import (
	"context"
	"fmt"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/expression"
	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/streaming/stream/input"
)

type InsertOperator struct {
	ctx     context.Context
	exprCtx expression.EvalContext
	insert  *plan.InsertNode

	child   operator.Operator
	inbound *operator.Queue
}

func NewInsertOperator(ctx context.Context, insert *plan.InsertNode, child operator.Operator) operator.Operator {
	return &InsertOperator{
		ctx:     ctx,
		insert:  insert,
		child:   child,
		inbound: operator.NewQueue(make(chan *types.Page, 256)),
	}
}

func (op *InsertOperator) Run(ctx context.Context, output chan<- *types.Page) {
	streamName := op.insert.Table.Name.Name
	// FIXME: get app from context
	inputHandle := input.GetManager().GetInputHandler("test", streamName)
	for {
		page, ok := op.inbound.Consume(ctx)
		fmt.Printf("insert...%v=>%v", ok, page)
		if !ok {
			break
		}
		inputHandle.Send(page)
	}
}

func (op *InsertOperator) GetLayout() []*plan.Symbol {
	return op.insert.GetOutputSymbols()
}

func (op *InsertOperator) Children() []operator.Operator {
	return []operator.Operator{op.child}
}

func (op *InsertOperator) GetInbounds() []chan *types.Page {
	return []chan *types.Page{op.inbound.GetInbound()}
}

func (op *InsertOperator) String() string {
	return "InsertOperator"
}
