package processor

import (
	"github.com/lindb/lindb/streaming/query/context"
)

type OutputProcessor struct {
	ctx *context.QueryContext

	child Processor
}

func NewOutputProcessor(ctx *context.QueryContext,
	child Processor,
) Processor {
	p := &OutputProcessor{
		ctx:   ctx,
		child: child,
	}
	return p
}

func (o *OutputProcessor) Run(output chan<- any) {
}

func (o *OutputProcessor) GetInbounds() []chan any {
	return nil
}

func (o *OutputProcessor) Children() []Processor {
	return []Processor{o.child}
}

func (o *OutputProcessor) String() string {
	return "Output"
}
