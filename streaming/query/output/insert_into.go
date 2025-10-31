package output

import (
	"fmt"

	"github.com/lindb/lindb/sql/planner/plan"
	"github.com/lindb/lindb/streaming/query/context"
	"github.com/lindb/lindb/streaming/query/processor"
)

type InsertInto struct {
	ctx  *context.QueryContext
	node *plan.InsertNode

	child processor.Processor

	inbound *processor.Queue
}

func NewInsertInto(ctx *context.QueryContext,
	node *plan.InsertNode,
	child processor.Processor,
) processor.Processor {
	return &InsertInto{
		ctx:     ctx,
		node:    node,
		child:   child,
		inbound: processor.NewQueue(make(chan any)),
	}
}

func (i *InsertInto) Run(output chan<- any) {
	streamName := i.node.Table.Name.Name
	handle := i.ctx.InputManager.GetInputHandler(streamName)
	for {
		source, ok := i.inbound.Consume(i.ctx.Context)
		if !ok {
			break
		}
		fmt.Printf("current:%v, insert into:%v\n", i.ctx.Query, streamName)
		handle.Send(source)
	}
}

func (i *InsertInto) GetInbounds() []chan any {
	return []chan any{i.inbound.GetInbound()}
}

func (i *InsertInto) Children() []processor.Processor {
	return []processor.Processor{i.child}
}

func (i *InsertInto) String() string {
	return fmt.Sprintf("InsertInto[stream=%s]", i.node.Table.Name.Name)
}
