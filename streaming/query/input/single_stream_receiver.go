package input

import (
	"fmt"

	"github.com/lindb/lindb/streaming/query/context"
	"github.com/lindb/lindb/streaming/query/processor"
)

type SingleStreamReceiver struct {
	ctx    *context.QueryContext
	stream string

	inbound *processor.Queue
}

func NewSingleStreamReceiver(ctx *context.QueryContext, stream string) *SingleStreamReceiver {
	return &SingleStreamReceiver{
		ctx:     ctx,
		stream:  stream,
		inbound: processor.NewQueue(make(chan any)),
	}
}

func (r *SingleStreamReceiver) Receive(event any) {
	fmt.Printf("single stream receiver, current:%v, receive event:%v\n", r.ctx.Query, event)
	r.inbound.Produce(event)
}

func (r *SingleStreamReceiver) Run(output chan<- any) {
	for {
		source, ok := r.inbound.Consume(r.ctx.Context)
		if !ok {
			break
		}
		output <- source
	}
}

func (r *SingleStreamReceiver) GetInbounds() []chan any {
	return []chan any{r.inbound.GetInbound()}
}

func (r *SingleStreamReceiver) Children() []processor.Processor {
	return nil
}

func (r *SingleStreamReceiver) String() string {
	return fmt.Sprintf("SingleStreamReceiver[stream=%s]", r.stream)
}
