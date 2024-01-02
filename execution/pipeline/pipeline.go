package pipeline

import (
	"fmt"
	"reflect"

	"github.com/lindb/lindb/spi"
	"github.com/lindb/lindb/sql/context"
)

type Pipeline struct {
	taskCtx     *context.TaskContext
	splitSource spi.SplitSource
	driverFct   *DriverFactory
}

func NewPipeline(taskCtx *context.TaskContext, splitSource spi.SplitSource, driverFct *DriverFactory) *Pipeline {
	return &Pipeline{
		taskCtx:     taskCtx,
		splitSource: splitSource,
		driverFct:   driverFct,
	}
}

func (p *Pipeline) Run() {
	if p.splitSource == nil {
		// source from exchange(local/remote)
		driver := p.driverFct.CreateDriver()
		sourceOperator := driver.GetSourceOperator()
		fmt.Printf("run diver====%v,%s\n", sourceOperator, reflect.TypeOf(sourceOperator))
		if sourceOperator != nil {
			// if driver has source operator, register it
			DriverManager.RegisterSourceOperator(p.taskCtx.TaskID, sourceOperator)
		}

		driver.Process()
	} else {
		// source from storage
		p.splitSource.Prepare()
		driver := p.driverFct.CreateDriver()

		for p.splitSource.HasSplit() {
			split := p.splitSource.GetNextSplit()

			if split != nil {
				driver.AddSplit(split)
				driver.Process()
			}
		}
	}
}
