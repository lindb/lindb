package query

import (
	"sync/atomic"

	"github.com/lindb/lindb/series"
)

type storageExecutorContext struct {
	resultCh chan *series.TimeSeriesEvent

	taskCounter int32 // pending task ref counter
}

func newStorageExecutorContext(resultCh chan *series.TimeSeriesEvent) *storageExecutorContext {
	return &storageExecutorContext{resultCh: resultCh}
}

func (c *storageExecutorContext) retainTask(tasks int32) {
	atomic.AddInt32(&c.taskCounter, tasks)
}

func (c *storageExecutorContext) completeTask() {
	newVal := atomic.AddInt32(&c.taskCounter, -1)
	// if all tasks completed, close result channel
	if newVal == 0 {
		close(c.resultCh)
	}
}
