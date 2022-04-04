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

package concurrent

import (
	"context"
	"sync"
	"time"

	"github.com/lindb/lindb/internal/linmetric"
	errorpkg "github.com/lindb/lindb/pkg/error"
	"github.com/lindb/lindb/pkg/logger"

	"go.uber.org/atomic"
)

//go:generate mockgen -source=./pool.go -destination=./pool_mock.go -package=concurrent

const (
	// size of the queue that workers register their availability to the dispatcher.
	readyWorkerQueueSize = 32
	// size of the tasks queue
	tasksCapacity = 8
	// sleeps in this interval when there are no available workers
	sleepInterval = time.Millisecond * 5
)

// Task represents a task function to be executed by a worker(goroutine).
type Task struct {
	// handle executes task function.
	handle func()
	// panicHandle executes callback if task happens panic.
	panicHandle func(err error)

	createTime time.Time
}

// NewTask creates a task.
func NewTask(handle func(), panicHandle func(err error)) *Task {
	return &Task{
		handle:      handle,
		panicHandle: panicHandle,
		createTime:  time.Now(),
	}
}

func (t *Task) Exec() {
	t.handle()
}

// Pool represents the goroutine pool that executes submitted tasks.
type Pool interface {
	// Submit enqueues a callable task for a worker to execute.
	//
	// Each submitted task is immediately given to a ready worker.
	// If there are no available workers, the dispatcher starts a new worker,
	// until the maximum number of workers are added.
	//
	// After the maximum number of workers are running, and no workers are ready,
	// execute function will be blocked.
	Submit(ctx context.Context, task *Task)
	// Stopped returns true if this pool has been stopped.
	Stopped() bool
	// Stop stops all goroutines gracefully,
	// all pending tasks will be finished before exit
	Stop()
}

// workerPool is a pool for goroutines.
type workerPool struct {
	name                string
	maxWorkers          int
	tasks               chan *Task              // tasks channel
	readyWorkers        chan *worker            // available worker
	idleTimeout         time.Duration           // idle goroutine recycle time
	onDispatcherStopped chan struct{}           // signal that dispatcher is stopped
	stopped             atomic.Bool             // mark if the pool is closed or not
	workersAlive        *linmetric.BoundGauge   // current workers count in use
	workersCreated      *linmetric.BoundCounter // workers created count since start
	workersKilled       *linmetric.BoundCounter // workers killed since start
	tasksConsumed       *linmetric.BoundCounter // tasks consumed count
	tasksRejected       *linmetric.BoundCounter // tasks rejected count
	tasksPanic          *linmetric.BoundCounter // tasks execute panic count
	tasksWaitingTime    *linmetric.BoundCounter // tasks waiting total time
	tasksExecutingTime  *linmetric.BoundCounter // tasks executing total time with waiting period
	ctx                 context.Context
	cancel              context.CancelFunc

	logger *logger.Logger
}

// NewPool returns a new worker pool,
// maxWorkers parameter specifies the maximum number workers that will execute tasks concurrently.
func NewPool(name string, maxWorkers int, idleTimeout time.Duration, scope linmetric.Scope) Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	pool := &workerPool{
		name:                name,
		maxWorkers:          maxWorkers,
		tasks:               make(chan *Task, tasksCapacity),
		readyWorkers:        make(chan *worker, readyWorkerQueueSize),
		idleTimeout:         idleTimeout,
		onDispatcherStopped: make(chan struct{}),
		stopped:             *atomic.NewBool(false),
		workersAlive:        scope.NewGauge("workers_alive"),
		workersCreated:      scope.NewCounter("workers_created"),
		workersKilled:       scope.NewCounter("workers_killed"),
		tasksConsumed:       scope.NewCounter("tasks_consumed"),
		tasksRejected:       scope.NewCounter("tasks_rejected"),
		tasksPanic:          scope.NewCounter("tasks_panic"),
		tasksWaitingTime:    scope.NewCounter("tasks_waiting_duration_sum"),
		tasksExecutingTime:  scope.NewCounter("tasks_executing_duration_sum"),
		ctx:                 ctx,
		cancel:              cancel,
		logger:              logger.GetLogger("Pool", name),
	}
	go pool.dispatch()
	return pool
}

func (p *workerPool) Submit(ctx context.Context, task *Task) {
	if task.handle == nil || p.Stopped() {
		return
	}
	select {
	case <-ctx.Done():
		p.tasksRejected.Incr()
		return
	case p.tasks <- task:
	}
}

// mustGetWorker makes sure that a ready worker is return
func (p *workerPool) mustGetWorker() *worker {
	var worker *worker
	for {
		select {
		// got a worker
		case worker = <-p.readyWorkers:
			return worker
		default:
			if int(p.workersAlive.Get()) >= p.maxWorkers {
				// no available workers
				time.Sleep(sleepInterval)
				continue
			}
			w := newWorker(p)
			return w
		}
	}
}

func (p *workerPool) dispatch() {
	defer func() {
		p.onDispatcherStopped <- struct{}{}
	}()

	idleTimeoutTimer := time.NewTimer(p.idleTimeout)
	defer idleTimeoutTimer.Stop()
	var (
		worker *worker
		task   *Task
	)

	for {
		idleTimeoutTimer.Reset(p.idleTimeout)
		select {
		case <-p.ctx.Done():
			return
		case task = <-p.tasks:
			worker = p.mustGetWorker()
			worker.execute(task)
		case <-idleTimeoutTimer.C:
			// timed out waiting, kill a ready worker
			if p.workersAlive.Get() > 0 {
				select {
				case worker = <-p.readyWorkers:
					worker.stop(func() {})
				default:
					// workers are busy now
				}
			}
		}
	}
}

func (p *workerPool) Stopped() bool {
	return p.stopped.Load()
}

// stopWorkers stops all workers
func (p *workerPool) stopWorkers() {
	var wg sync.WaitGroup
	for p.workersAlive.Get() > 0 {
		wg.Add(1)
		worker := <-p.readyWorkers
		worker.stop(func() {
			wg.Done()
		})
	}
	wg.Wait()
}

// consumedRemainingTasks consumes all buffered tasks in the channel
func (p *workerPool) consumedRemainingTasks() {
	for {
		select {
		case task := <-p.tasks:
			p.execTask(task)
		default:
			return
		}
	}
}

func (p *workerPool) execTask(task *Task) {
	defer func() {
		var err error
		r := recover()
		if r != nil {
			p.tasksPanic.Incr()
			err = errorpkg.Error(r)
			p.logger.Error("panic when execute task",
				logger.Error(err), logger.Stack())
			if task.panicHandle != nil {
				task.panicHandle(err)
			}
		}
	}()
	p.tasksWaitingTime.Add(float64(time.Since(task.createTime).Nanoseconds() / 1e6))
	task.Exec()
	p.tasksExecutingTime.Add(float64(time.Since(task.createTime).Nanoseconds() / 1e6))

	p.tasksConsumed.Incr()
}

// Stop tells the dispatcher to exit with pending tasks done.
func (p *workerPool) Stop() {
	if p.stopped.Swap(true) {
		return
	}
	// close dispatcher
	p.cancel()
	// wait dispatcher's exit
	<-p.onDispatcherStopped
	// close all workers
	p.stopWorkers()
	// consume remaining tasks
	p.consumedRemainingTasks()
}

// worker represents the worker that executes the task
type worker struct {
	pool   *workerPool
	tasks  chan *Task
	stopCh chan struct{}
}

// newWorker creates the worker that executes tasks given by the dispatcher
// When a new worker starts, it registers itself on the createdWorkers channel.
func newWorker(pool *workerPool) *worker {
	w := &worker{
		pool:   pool,
		tasks:  make(chan *Task),
		stopCh: make(chan struct{}),
	}
	w.pool.workersAlive.Incr()
	w.pool.workersCreated.Incr()
	go w.process()
	return w
}

// execute submits the task to queue
func (w *worker) execute(task *Task) {
	w.tasks <- task
}

func (w *worker) stop(callable func()) {
	defer callable()
	w.stopCh <- struct{}{}
	w.pool.workersKilled.Incr()
	w.pool.workersAlive.Decr()
}

// process task from queue
func (w *worker) process() {
	var task *Task
	for {
		select {
		case <-w.stopCh:
			return
		case task = <-w.tasks:
			w.pool.execTask(task)
			// register worker-self to readyWorkers again
			w.pool.readyWorkers <- w
		}
	}
}
