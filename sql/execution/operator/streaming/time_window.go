package streaming

import (
	"context"
	"time"

	"github.com/lindb/lindb/spi/types"
	"github.com/lindb/lindb/sql/execution/operator"
	"github.com/lindb/lindb/sql/execution/operator/streaming/executor"
	"github.com/lindb/lindb/sql/planner/plan"
)

type TimeWindowOperator struct {
	node  plan.PlanNode
	child operator.Operator

	executor executor.Executor

	inbound *operator.Queue

	ticker *time.Ticker
}

func NewTimeWindowOperator(node plan.PlanNode, executor executor.Executor, child operator.Operator) operator.Operator {
	ticker := time.NewTicker(time.Second * 5)
	return &TimeWindowOperator{
		node:     node,
		executor: executor,
		child:    child,
		inbound:  operator.NewQueue(make(chan *types.Page, 1024)),
		ticker:   ticker,
	}
}

func (t *TimeWindowOperator) Run(ctx context.Context, output chan<- *types.Page) {
	defer t.ticker.Stop()

	for {
		select {
		case source := <-t.inbound.GetInbound():
			if source == nil {
				return
			}
			t.executor.Enter(source)
			// fmt.Println("time window...")
		case <-t.ticker.C:
			t.executor.Leave(output)
			// fmt.Println("get executor result")
		}
	}
}

func (t *TimeWindowOperator) GetLayout() []*plan.Symbol {
	return t.node.GetOutputSymbols()
}

func (t *TimeWindowOperator) Children() []operator.Operator {
	return []operator.Operator{t.child}
}

func (t *TimeWindowOperator) GetInbounds() []chan *types.Page {
	return []chan *types.Page{t.inbound.GetInbound()}
}

func (t *TimeWindowOperator) String() string {
	return "TimeWindowOperator"
}
