package task

import (
	"context"
	"fmt"

	"go.uber.org/zap"

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
	log := logger.GetLogger()
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
				log.Error("watch events", zap.Error(evt.Err))
			}
		}
	}
}

func (e *Executor) dispatch(kvevt state.EventKeyValue) {
	log := logger.GetLogger()
	var task Task
	(&task).UnsafeUnmarshal(kvevt.Value)
	if task.State > StateRunning {
		log.Debug("stale task", zap.String("name", kvevt.Key))
		return
	}
	proc, ok := e.processors[task.Kind]
	if !ok {
		log.Warn("processor not found", zap.String("kind", string(task.Kind)))
		return
	}
	// Each task processor has an infinite events queue, so it won't block others.
	err := proc.Submit(taskEvent{key: kvevt.Key, task: task, rev: kvevt.Rev})
	if err != nil {
		log.Warn("dispatch task", zap.Error(err))
	} else {
		log.Info("dispatch task", zap.String("name", kvevt.Key))
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
