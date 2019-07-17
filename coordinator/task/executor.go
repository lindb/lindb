package task

import (
	"context"
	"fmt"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"
)

// Executor executes tasks on node.
type Executor struct {
	keypfx     string
	cli        state.Repository
	node       *models.Node
	processors map[Kind]*taskProcessor

	ctx    context.Context
	cancel context.CancelFunc

	log *logger.Logger
}

// NewExecutor creates a new Executor, the keypfx must be the same as Controller's.
func NewExecutor(ctx context.Context, node *models.Node, cli state.Repository) *Executor {
	ctx, cancel := context.WithCancel(ctx)
	return &Executor{
		keypfx:     fmt.Sprintf("/task-coordinator/%s/executor/%s/", version, node.String()),
		cli:        cli,
		node:       node,
		processors: map[Kind]*taskProcessor{},
		ctx:        ctx,
		log:        logger.GetLogger("coordinator/task/executor"),
		cancel:     cancel,
	}
}

// Register must be called before Run.
func (e *Executor) Register(procs ...Processor) {
	for _, proc := range procs {
		e.processors[proc.Kind()] = newTaskProcessor(e.ctx, proc, e.cli)
	}
}

// Run must be called after Register, otherwise it may panic, O(∩_∩)O~.
func (e *Executor) Run() {
	evtc := e.cli.WatchPrefix(e.ctx, e.keypfx)
	for {
		select {
		case <-e.ctx.Done():
			return
		case evt := <-evtc:
			if evt == nil { // Context canceled
				return
			}
			if evt.Err == nil {
				switch evt.Type {
				case state.EventTypeAll:
					fallthrough
				case state.EventTypeModify:
					for _, kv := range evt.KeyValues {
						e.dispatch(kv)
					}
				case state.EventTypeDelete:
				}
			} else {
				e.log.Error("watch events", logger.Error(evt.Err))
			}
		}
	}
}

func (e *Executor) dispatch(kvevt state.EventKeyValue) {
	var task Task
	(&task).UnsafeUnmarshal(kvevt.Value)
	if task.State > StateRunning {
		e.log.Debug("stale task", logger.String("name", kvevt.Key))
		return
	}
	proc, ok := e.processors[task.Kind]
	if !ok {
		e.log.Warn("processor not found", logger.String("kind", string(task.Kind)))
		return
	}
	// Each task processor has an infinite events queue, so it won't block others.
	err := proc.Submit(taskEvent{key: kvevt.Key, task: task, rev: kvevt.Rev})
	if err != nil {
		e.log.Warn("dispatch task", logger.Error(err))
	} else {
		e.log.Info("dispatch task", logger.String("name", kvevt.Key))
	}
}

// Close closes Executor.
func (e *Executor) Close() error {
	e.cancel()
	for _, p := range e.processors {
		p.Stop()
	}
	return nil
}
