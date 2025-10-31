package processor

import (
	"fmt"
	"time"

	"github.com/lindb/lindb/streaming/executor"
	"github.com/lindb/lindb/streaming/query/context"
)

type TimeWindowProcessor struct {
	ctx *context.QueryContext

	executor executor.Executor
	child    Processor
	ticker   *time.Ticker

	inbound *Queue
}

func NewTimeWindowProcessor(
	ctx *context.QueryContext,
	executor executor.Executor,
	child Processor,
) Processor {
	p := &TimeWindowProcessor{
		ctx:      ctx,
		executor: executor,
		child:    child,
		inbound:  NewQueue(make(chan any)),
	}
	p.schedule()
	return p
}

func (t *TimeWindowProcessor) schedule() {
	t.ticker = time.NewTicker(time.Second * 5)
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			t.executor.ResultSet(func(result any) {
	// 				// t.Next(result)
	// 			})
	// 		case <-t.ctx.Context.Done():
	// 			ticker.Stop()
	// 			return
	// 		}
	// 	}
	// }()
}

// Process implements Processor.
//
//	func (t *TimeWindowProcessor) Process(event any) {
//		fmt.Printf("current:%v, time window event=%v\n", t.ctx.Query, event)
//		t.executor.Process(event)
//		// t.Next(event)
//	}
func (t *TimeWindowProcessor) Run(output chan<- any) {
	for {
		select {
		case <-t.ticker.C:
			t.executor.ResultSet(func(result any) {
				output <- result
				// t.Next(result)
			})
		case event := <-t.inbound.GetInbound():
			fmt.Printf("current:%v, time window event=%v\n", t.ctx.Query, event)
			if event != nil {
				t.executor.Process(event)
			}
		case <-t.ctx.Context.Done():
			t.ticker.Stop()
			return
		}
	}
}

func (t *TimeWindowProcessor) GetInbounds() []chan any {
	return []chan any{t.inbound.GetInbound()}
}

func (t *TimeWindowProcessor) Children() []Processor {
	return []Processor{t.child}
}

func (t *TimeWindowProcessor) String() string {
	return "TimeWindow"
}
