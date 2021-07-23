// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package task

import (
	"context"
	"fmt"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

// Executor watches task event and dispatches the task event to target task processor
type Executor struct {
	keyPrefix  string
	repo       state.Repository
	node       *models.StatefulNode
	processors map[Kind]*taskProcessor

	ctx    context.Context
	cancel context.CancelFunc

	log *logger.Logger
}

// NewExecutor creates a new Executor, the task key prefix must be the same as Controller's.
func NewExecutor(ctx context.Context, node *models.StatefulNode, repo state.Repository) *Executor {
	ctx, cancel := context.WithCancel(ctx)
	return &Executor{
		keyPrefix:  fmt.Sprintf("%s/executor/%s/", taskCoordinatorKey, node.Indicator()),
		repo:       repo,
		node:       node,
		processors: map[Kind]*taskProcessor{},
		ctx:        ctx,
		cancel:     cancel,
		log:        logger.GetLogger("coordinator", "TaskExecutor"),
	}
}

// Register registers task processor.
// Notice: must be called before Run.
func (e *Executor) Register(procs ...Processor) {
	for _, proc := range procs {
		e.processors[proc.Kind()] = newTaskProcessor(e.ctx, proc, e.repo)
	}
}

// Run must be called after Register, otherwise it may panic, O(∩_∩)O~.
func (e *Executor) Run() {
	eventCh := e.repo.WatchPrefix(e.ctx, e.keyPrefix, true)
	for {
		select {
		case <-e.ctx.Done(): // Context canceled
			return
		case evt := <-eventCh:
			if evt == nil { // chan closed
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
					// delete not used
				}
			} else {
				e.log.Error("watch task events", logger.Error(evt.Err))
			}
		}
	}
}

// dispatch dispatches task event to target task processor
func (e *Executor) dispatch(eventKV state.EventKeyValue) {
	task := Task{}
	if err := encoding.JSONUnmarshal(eventKV.Value, &task); err != nil {
		e.log.Error("unmarshal task data", logger.Any("data", eventKV.Value))
		return
	}

	if task.State > StateRunning {
		e.log.Debug("stale task", logger.String("name", eventKV.Key))
		return
	}
	proc, ok := e.processors[task.Kind]
	if !ok {
		e.log.Warn("processor not found", logger.String("kind", string(task.Kind)))
		return
	}
	// Each task processor has an infinite events queue, so it won't block others.
	err := proc.Submit(taskEvent{key: eventKV.Key, task: task, rev: eventKV.Rev})
	if err != nil {
		e.log.Warn("dispatch task", logger.Error(err))
	} else {
		e.log.Info("dispatch task successfully", logger.String("name", eventKV.Key))
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
