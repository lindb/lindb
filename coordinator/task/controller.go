package task

import (
	"context"
	"fmt"
	"strings"
	"sync/atomic"

	etcdcliv3 "github.com/coreos/etcd/clientv3"
	"github.com/damnever/goctl/zerocopy"
	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"
)

const (
	version = "v1" // TODO
	// FIXME(damnever): magic number, see also: --max-txn-ops in etcd
	maxTasksLimit = 127
)

var (
	ErrControllerClosed       = fmt.Errorf("coordinator/task: controller closed")
	ErrMaxTasksLimitExceeded  = fmt.Errorf("coordinator/task: tasks number can not greater than %d", maxTasksLimit)
	ErrTaskNameAlreadyExisted = fmt.Errorf("coordinator/task: task name already existed")
)

// ToBytes can convert itself into bytes.
type ToBytes interface {
	Bytes() []byte
}

// ControllerTaskParam is the param for tasks, the node id is the unique id for node.
type ControllerTaskParam struct {
	NodeID string
	Params ToBytes
}

// Controller is responsible for submitting tasks, noticing responding when task status changes.
//
// TODO(damnever): API to notify task status changes.
//  - we can simply watch the key: /<keypfx>/task-coordinator/<version>/status/kinds/<task-kind>/names/<task-name>
type Controller struct {
	keypfx string
	cli    state.Repository

	ctx    context.Context
	cancel context.CancelFunc
	donec  chan struct{}
	closed int32
}

// NewController creates a new controller.
func NewController(ctx context.Context, keypfx string, cli state.Repository) *Controller {
	ctx, cancel := context.WithCancel(ctx)
	c := &Controller{
		keypfx: fmt.Sprintf("%s/task-coordinator/%s", keypfx, version),
		cli:    cli,
		ctx:    ctx,
		cancel: cancel,
		donec:  make(chan struct{}),
	}
	go c.run()
	return c
}

// Submit submits a task with params and node ids, a readable name with
// context information is recommended.
func (c *Controller) Submit(kind Kind, name string, params []ControllerTaskParam) error {
	if atomic.LoadInt32(&c.closed) == 1 {
		return ErrControllerClosed
	}
	if len(params) >= maxTasksLimit {
		return ErrMaxTasksLimitExceeded
	}
	if len(params) == 0 {
		return nil
	}

	// TODO(damnever): kinds validation
	grp := groupedTasks{State: StateRunning}
	ops := []etcdcliv3.Op{{}}
	for _, param := range params {
		task := Task{
			Kind:     kind,
			Name:     name,
			Executor: param.NodeID,
			Params:   param.Params.Bytes(),
			State:    StateCreated,
		}
		grp.Tasks = append(grp.Tasks, task)
		ops = append(ops, etcdcliv3.OpPut(
			c.taskKey(kind, name, param.NodeID),
			zerocopy.UnsafeBtoa(task.UnsafeMarshal()),
		))
	}
	key := c.statusKey(kind, name)
	ops[0] = etcdcliv3.OpPut(
		key,
		zerocopy.UnsafeBtoa(grp.UnsafeMarshal()),
	)

	resp, err := c.cli.Txn(c.ctx).If(
		etcdcliv3.Compare(etcdcliv3.CreateRevision(key), "=", 0),
	).Then(ops...).Commit()
	if err != nil {
		return err
	}
	if !resp.Succeeded {
		return ErrTaskNameAlreadyExisted
	}
	return nil
}

// Close shutdown Controller.
func (c *Controller) Close() error {
	if atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		c.cancel()
		<-c.donec
	}
	return nil
}

func (c *Controller) run() {
	close(c.donec)

	evtc := c.cli.WatchPrefix(c.ctx, c.keypfx)
	waiters := newWaiters()
	log := logger.GetLogger()
	for {
		select {
		case <-c.ctx.Done():
			return
		case evt := <-evtc:
			if evt == nil {
				log.Warn("task event channel closed")
				return
			}
			if evt.Err == nil {
				switch evt.Type {
				case state.EventTypeAll:
					waiters = newWaiters()
					fallthrough
				case state.EventTypeModify:
					for _, kv := range evt.KeyValues {
						// FIXME(damnever): more robust checking
						if strings.Contains(kv.Key, "status") {
							var tasks groupedTasks
							(&tasks).UnsafeUnmarshal(kv.Value)
							waiters.TryAdd(kv.Key, tasks, kv.Rev)
						} else {
							var task Task
							(&task).UnsafeUnmarshal(kv.Value)
							if err := waiters.TryNotify(c, task); err != nil {
								log.Error("update status", zap.Error(err))
							}
						}
					}
				case state.EventTypeDelete:
					// NOTE(damnever): delete events are useless
				}
			} else {
				log.Warn("error event", zap.Error(evt.Err))
			}
		}
	}
}

func (c *Controller) statusKey(kind Kind, name string) string {
	return fmt.Sprintf("%s/status/kinds/%s/names/%s", c.keypfx, kind, name)
}

func (c *Controller) taskKey(kind Kind, name, nodeID string) string {
	return fmt.Sprintf("%s/executor/%s/kinds/%s/names/%s", c.keypfx, nodeID, kind, name)
}

type statusWaiter struct {
	kind    Kind
	name    string
	key     string
	rev     int64
	tasks   groupedTasks
	waiting map[string]struct{}
}

func newStatusWaiter(key string, tasks groupedTasks, rev int64) *statusWaiter {
	w := &statusWaiter{
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

func (w *statusWaiter) Confirm(task Task) {
	if task.Kind != w.kind || task.Name != w.name {
		return
	}
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

func (w *statusWaiter) UpdateStatus(c *Controller) error {
	ops := []etcdcliv3.Op{{}}
	w.tasks.State = StateDoneOK
	for _, task := range w.tasks.Tasks {
		ops = append(ops, etcdcliv3.OpDelete(c.taskKey(task.Kind, task.Name, task.Executor)))
		if task.ErrMsg != "" {
			w.tasks.State = StateDoneErr
		}
	}
	ops[0] = etcdcliv3.OpPut(w.key, zerocopy.UnsafeBtoa(w.tasks.UnsafeMarshal()))
	resp, err := c.cli.Txn(c.ctx).If(
		etcdcliv3.Compare(etcdcliv3.ModRevision(w.key), "=", w.rev),
	).Then(ops...).Commit()
	if err != nil {
		return err
	}
	if !resp.Succeeded {
		return fmt.Errorf("coordinator/task: someone changed the %s", w.key)
	}
	return nil
}

func (w *statusWaiter) IsDone() bool {
	return len(w.waiting) == 0
}

type kindStatusWaiter map[string]*statusWaiter

type waiters map[Kind]kindStatusWaiter

func newWaiters() waiters {
	return map[Kind]kindStatusWaiter{}
}

func (w waiters) TryAdd(key string, tasks groupedTasks, rev int64) {
	if tasks.State > StateRunning { // NOTE(damnever): Ignore newly finished tasks.
		return
	}

	kind := tasks.Tasks[0].Kind
	kw, ok := w[kind]
	if !ok {
		kw = kindStatusWaiter{}
		w[kind] = kw
	}
	// FIXME(damnever): check duplicate
	kw[tasks.Tasks[0].Name] = newStatusWaiter(key, tasks, rev)
}

func (w waiters) TryNotify(c *Controller, task Task) (err error) {
	if task.State <= StateRunning { // NOTE(damnever): Ignore newly created tasks.
		return
	}
	kw, ok := w[task.Kind]
	if !ok {
		return
	}
	sw, ok := kw[task.Name]
	if !ok {
		return
	}

	sw.Confirm(task)
	if sw.IsDone() {
		if err = sw.UpdateStatus(c); err == nil {
			delete(kw, task.Name)
		}
	}
	return
}
