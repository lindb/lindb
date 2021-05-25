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
	"strings"
	"sync/atomic"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./controller.go -destination=./controller_mock.go -package=task

var (
	// ErrControllerClosed causes error when does some ops after controller is closed
	ErrControllerClosed = fmt.Errorf("coordinator/task: controller closed")
	// ErrMaxTasksLimitExceeded causes error when exceeds the task limit
	ErrMaxTasksLimitExceeded = fmt.Errorf("coordinator/task: tasks number can not greater than %d", maxTasksLimit)
)

var log = logger.GetLogger("coordinator", "TaskController")

// ControllerFactory represents a task controller create factory
type ControllerFactory interface {
	// CreateController creates a new task controller
	CreateController(ctx context.Context, cli state.Repository) Controller
}

// controllerFactory implements the interface
type controllerFactory struct {
}

// NewControllerFactory creates a task controller create factory
func NewControllerFactory() ControllerFactory {
	return &controllerFactory{}
}

// NewController creates a new controller.
func (f *controllerFactory) CreateController(ctx context.Context, repo state.Repository) Controller {
	ctx, cancel := context.WithCancel(ctx)
	c := &controller{
		keyPrefix:    taskCoordinatorKey,
		statusPrefix: fmt.Sprintf("%s/status/kinds", taskCoordinatorKey),
		repo:         repo,
		ctx:          ctx,
		cancel:       cancel,
		donec:        make(chan struct{}),
	}
	go c.run()
	return c
}

// ToBytes can convert itself into bytes.
type ToBytes interface {
	Bytes() []byte
}

// ControllerTaskParam is the param for tasks, the node id is the unique id for node.
type ControllerTaskParam struct {
	NodeID string
	Params ToBytes
}

// Controller is responsible for submitting tasks
type Controller interface {
	// Submit submits a task with params
	Submit(kind Kind, name string, params []ControllerTaskParam) error
	// Close closes controller, then releases the resource
	Close() error
	// taskKey returns the key of task
	taskKey(kind Kind, name, nodeID string) string
	// statusKey returns the key of task status
	statusKey(kind Kind, name string) string
}

// controller is responsible for submitting tasks, noticing responding when task status changes.
//
// API to notify task status changes.
//  - we can simply watch the key: /task-coordinator/<version>/status/kinds/<task-kind>/names/<task-name>
type controller struct {
	keyPrefix    string
	statusPrefix string
	repo         state.Repository

	ctx    context.Context
	cancel context.CancelFunc
	donec  chan struct{}
	closed int32
}

// Submit submits a task with params and node ids, a readable name with
// context information is recommended.
func (c *controller) Submit(kind Kind, name string, params []ControllerTaskParam) error {
	if atomic.LoadInt32(&c.closed) == 1 {
		return ErrControllerClosed
	}
	if len(params) >= maxTasksLimit {
		return ErrMaxTasksLimitExceeded
	}
	if len(params) == 0 {
		return nil
	}

	txn := c.repo.NewTransaction()
	// TODO(damnever): kinds validation
	grp := groupedTasks{State: StateRunning}
	for _, param := range params {
		task := Task{
			Kind:     kind,
			Name:     name,
			Executor: param.NodeID,
			Params:   param.Params.Bytes(),
			State:    StateCreated,
		}
		grp.Tasks = append(grp.Tasks, task)
	}
	// first, must put task status track
	txn.Put(c.statusKey(kind, name), encoding.JSONMarshal(&grp))
	// then, add tasks
	for i := range grp.Tasks {
		task := grp.Tasks[i]
		txn.Put(c.taskKey(kind, name, task.Executor), encoding.JSONMarshal(&task))
	}
	// finally, commit txn
	err := c.repo.Commit(c.ctx, txn)
	if err != nil {
		return err
	}
	return nil
}

// Close shutdowns task controller
func (c *controller) Close() error {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		log.Info("closing task controller")
		c.cancel()
		<-c.donec
	}
	return nil
}

func (c *controller) run() {
	close(c.donec)

	log.Info("task controller running")

	defer log.Info("task controller loop exit")

	// watch "/task-coordinator/:version" for listening task event change
	eventCh := c.repo.WatchPrefix(c.ctx, c.keyPrefix, true)
	waiters := newWaiter(c.ctx, c.repo)
	for {
		select {
		case <-c.ctx.Done():
			return
		case evt := <-eventCh:
			if evt == nil {
				log.Warn("task event channel closed")
				return
			}
			if evt.Err == nil {
				switch evt.Type {
				case state.EventTypeAll:
					fallthrough
				case state.EventTypeModify:
					for _, kv := range evt.KeyValues {
						if strings.HasPrefix(kv.Key, c.statusPrefix) {
							tasks := groupedTasks{}
							if err := encoding.JSONUnmarshal(kv.Value, &tasks); err != nil {
								log.Error("unmarshal grouped tasks")
								continue
							}
							waiters.TryAdd(kv.Key, tasks, kv.Rev)
						} else {
							task := Task{}
							if err := encoding.JSONUnmarshal(kv.Value, &task); err != nil {
								log.Error("unmarshal task")
								continue
							}
							if err := waiters.TryNotify(c, task); err != nil {
								log.Error("update status", logger.Error(err))
							}
						}
					}
				case state.EventTypeDelete:
					// NOTE(damnever): delete events are useless
				}
			} else {
				log.Warn("got event with error", logger.Error(evt.Err))
			}
		}
	}
}

// statusKey returns the key of task status
func (c *controller) statusKey(kind Kind, name string) string {
	return fmt.Sprintf("%s/status/kinds/%s/names/%s", c.keyPrefix, kind, name)
}

// taskKey returns the key of task
func (c *controller) taskKey(kind Kind, name, nodeID string) string {
	return fmt.Sprintf("%s/executor/%s/kinds/%s/names/%s", c.keyPrefix, nodeID, kind, name)
}

type statusWaiter struct {
	ctx     context.Context
	repo    state.Repository
	kind    Kind
	name    string
	key     string
	rev     int64
	tasks   groupedTasks
	waiting map[string]struct{}
}

func newStatusWaiter(ctx context.Context,
	key string, tasks groupedTasks, rev int64,
	repo state.Repository) *statusWaiter {
	w := &statusWaiter{
		ctx:     ctx,
		repo:    repo,
		kind:    tasks.Tasks[0].Kind,
		name:    tasks.Tasks[0].Name,
		key:     key,
		rev:     rev,
		tasks:   tasks,
		waiting: map[string]struct{}{},
	}
	for _, task := range tasks.Tasks {
		w.waiting[task.Executor] = struct{}{}
	}
	return w
}

func (w *statusWaiter) confirm(task Task) {
	if _, ok := w.waiting[task.Executor]; !ok {
		return
	}
	delete(w.waiting, task.Executor)
	for i, t := range w.tasks.Tasks {
		if task.Executor == t.Executor {
			w.tasks.Tasks[i] = task
			break
		}
	}
}

func (w *statusWaiter) UpdateStatus(c Controller) error {
	txn := w.repo.NewTransaction()
	w.tasks.State = StateDoneOK
	for _, task := range w.tasks.Tasks {
		txn.Delete(c.taskKey(task.Kind, task.Name, task.Executor))
		if task.ErrMsg != "" {
			w.tasks.State = StateDoneErr
		}
	}
	txn.Put(w.key, encoding.JSONMarshal(&w.tasks))
	txn.ModRevisionCmp(w.key, "=", w.rev)

	if err := w.repo.Commit(w.ctx, txn); err != nil {
		return err
	}
	return nil
}

func (w *statusWaiter) IsDone() bool {
	return len(w.waiting) == 0
}

type kindStatusWaiter map[string]*statusWaiter

type waiter struct {
	ctx     context.Context
	repo    state.Repository
	waiters map[Kind]kindStatusWaiter
}

func newWaiter(ctx context.Context, repo state.Repository) *waiter {
	return &waiter{
		waiters: make(map[Kind]kindStatusWaiter),
		ctx:     ctx,
		repo:    repo,
	}
}

func (w *waiter) TryAdd(key string, tasks groupedTasks, rev int64) {
	if tasks.State > StateRunning { // NOTE(damnever): Ignore newly finished tasks.
		return
	}

	kind := tasks.Tasks[0].Kind
	kw, ok := w.waiters[kind]
	if !ok {
		kw = kindStatusWaiter{}
		w.waiters[kind] = kw
	}
	// FIXME(damnever): check duplicate
	kw[tasks.Tasks[0].Name] = newStatusWaiter(w.ctx, key, tasks, rev, w.repo)
}

func (w *waiter) TryNotify(c Controller, task Task) (err error) {
	if task.State <= StateRunning { // NOTE(damnever): Ignore newly created tasks.
		return
	}
	kw, ok := w.waiters[task.Kind]
	if !ok {
		return
	}
	sw, ok := kw[task.Name]
	if !ok {
		return
	}

	sw.confirm(task)
	if sw.IsDone() {
		if err = sw.UpdateStatus(c); err == nil {
			delete(kw, task.Name)
		}
	}
	return
}
