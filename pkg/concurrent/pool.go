package concurrent

import (
	"context"
	"sync"
	"time"

	"go.uber.org/atomic"
)

const (
	// size of the queue that workers register their availability to the dispatcher.
	readyWorkerQueueSize = 32
	// size of the tasks queue
	tasksCapacity = 8
	// sleeps in this interval when there are no available workers
	sleepInterval = time.Millisecond * 5
)

// Task represents a task function to be executed by a worker(goroutine).
type Task func()

// Pool represents the goroutine pool that executes submitted tasks.
type Pool interface {
	// Submit enqueues a callable task for a worker to execute.
	//
	// Each submitted task is immediately given to an ready worker.
	// If there are no available workers, the dispatcher starts a new worker,
	// until the maximum number of workers are added.
	//
	// After the maximum number of workers are running, and no workers are ready,
	// execute function will be blocked.
	Submit(task Task)
	// SubmitAndWait executes the task and waits for it to be executed.
	SubmitAndWait(task Task)
	// Stopped returns true if this pool has been stopped.
	Stopped() bool
	// Stop stops all goroutines gracefully,
	// all pending tasks will be finished before exit
	Stop()
	// Statistics returns the statistics data since started
	Statistics() *PoolStat
}

// workerPool is a pool for goroutines.
type workerPool struct {
	name                string
	maxWorkers          int
	tasks               chan Task     // tasks channel
	readyWorkers        chan *worker  // available worker
	idleTimeout         time.Duration // idle goroutine recycle time
	onDispatcherStopped chan struct{} // signal that dispatcher is stopped
	stopped             atomic.Bool   // mark if the pool is closed or not
	workersAlive        atomic.Int32  // current workers count in use
	workersCreated      atomic.Int32  // workers created count since start
	workersKilled       atomic.Int32  // workers killed since start
	tasksConsumed       atomic.Int32  // tasks consumed count
	ctx                 context.Context
	cancel              context.CancelFunc
}

// NewPool returns a new worker pool,
// maxWorkers parameter specifies the maximum number workers that will execute tasks concurrently.
func NewPool(name string, maxWorkers int, idleTimeout time.Duration) Pool {
	if maxWorkers < 1 {
		maxWorkers = 1
	}
	ctx, cancel := context.WithCancel(context.Background())
	pool := &workerPool{
		name:                name,
		maxWorkers:          maxWorkers,
		tasks:               make(chan Task, tasksCapacity),
		readyWorkers:        make(chan *worker, readyWorkerQueueSize),
		idleTimeout:         idleTimeout,
		onDispatcherStopped: make(chan struct{}),
		stopped:             *atomic.NewBool(false),
		workersAlive:        *atomic.NewInt32(0),
		workersCreated:      *atomic.NewInt32(0),
		workersKilled:       *atomic.NewInt32(0),
		tasksConsumed:       *atomic.NewInt32(0),
		ctx:                 ctx,
		cancel:              cancel,
	}
	go pool.dispatch()
	return pool
}

func (p *workerPool) Statistics() *PoolStat {
	return &PoolStat{
		AliveWorkers:   int(p.workersAlive.Load()),
		CreatedWorkers: int(p.workersCreated.Load()),
		KilledWorkers:  int(p.workersKilled.Load()),
		ConsumedTasks:  int(p.tasksConsumed.Load())}
}

func (p *workerPool) Submit(task Task) {
	if task == nil || p.Stopped() {
		return
	}
	p.tasks <- task
}

func (p *workerPool) SubmitAndWait(task Task) {
	if task == nil || p.Stopped() {
		return
	}
	worker := p.mustGetWorker()
	doneChan := make(chan struct{})
	worker.execute(func() {
		task()
		close(doneChan)
	})
	<-doneChan
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
			if int(p.workersAlive.Load()) >= p.maxWorkers {
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
		task   Task
	)

	for {
		idleTimeoutTimer.Reset(p.idleTimeout)
		select {
		case <-p.ctx.Done():
			return
		case task = <-p.tasks:
			worker := p.mustGetWorker()
			worker.execute(task)
		case <-idleTimeoutTimer.C:
			// timed out waiting, kill a ready worker
			if p.workersAlive.Load() > 0 {
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
	for p.workersAlive.Load() > 0 {
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
			task()
			p.tasksConsumed.Inc()
		default:
			return
		}
	}
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
	tasks  chan Task
	stopCh chan struct{}
}

// newWorker creates the worker that executes tasks given by the dispatcher
// When a new worker starts, it registers itself on the createdWorkers channel.
func newWorker(pool *workerPool) *worker {
	w := &worker{
		pool:   pool,
		tasks:  make(chan Task),
		stopCh: make(chan struct{}),
	}
	w.pool.workersAlive.Inc()
	w.pool.workersCreated.Inc()
	go w.process()
	return w
}

// execute submits the task to queue
func (w *worker) execute(task Task) {
	w.tasks <- task
}

func (w *worker) stop(callable func()) {
	defer callable()
	w.stopCh <- struct{}{}
	w.pool.workersKilled.Inc()
	w.pool.workersAlive.Dec()
}

// process process task from queue
func (w *worker) process() {
	var task Task
	for {
		select {
		case <-w.stopCh:
			return
		case task = <-w.tasks:
			task()
			w.pool.tasksConsumed.Inc()
			// register worker-self to readyWorkers again
			w.pool.readyWorkers <- w
		}
	}
}
