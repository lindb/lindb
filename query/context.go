package query

import (
	"github.com/lindb/lindb/series"

	"go.uber.org/atomic"
)

type storageExecutorContext struct {
	resultCh chan *series.TimeSeriesEvent

	taskCounter atomic.Int32 // pending task ref counter
}

func newStorageExecutorContext(resultCh chan *series.TimeSeriesEvent) *storageExecutorContext {
	return &storageExecutorContext{resultCh: resultCh}
}

func (c *storageExecutorContext) retainTask(tasks int32) {
	c.taskCounter.Add(tasks)
}

func (c *storageExecutorContext) completeTask() {
	newVal := c.taskCounter.Sub(1)
	// if all tasks completed, close result channel
	if newVal == 0 {
		close(c.resultCh)
	}
}
