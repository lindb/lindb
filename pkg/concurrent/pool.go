package concurrent

import (
	"go.uber.org/atomic"

	"github.com/lindb/lindb/pkg/logger"
)

var log = logger.GetLogger("concurrent", "pool")

// Task represents a task function to be executed by a worker(goroutine).
type Task func()

// Pool represents the goroutine pool that executes submitted tasks.
type Pool interface {
	// Name returns the name of pool
	Name() string
	// Execute executes the given task
	Execute(task Task)
	// Shutdown shutdowns all goroutines gracefully
	Shutdown()
}

// pool implements a simple goroutine pool
type pool struct {
	name string

	tasks   chan Task
	workers chan *worker
	stop    chan struct{}

	closing atomic.Bool

	// pool stats
	pending   atomic.Int64
	completed atomic.Int64

	worker atomic.Int64
}

// NewPool creates a pool
func NewPool(name string, nRoutines int, queueSize int) Pool {
	pool := &pool{
		name:    name,
		tasks:   make(chan Task, queueSize),
		workers: make(chan *worker, nRoutines),
		stop:    make(chan struct{}),
	}
	pool.worker.Add(int64(nRoutines))
	// init worker pool
	for i := 0; i < nRoutines; i++ {
		_ = newWorker(pool)
	}
	// start task dispatch routine
	go pool.dispatch()
	return pool
}

// Name returns the name
func (p *pool) Name() string {
	return p.name
}

// Execute puts the task to queue, if pool is closing reject the task
func (p *pool) Execute(task Task) {
	// pool is closing, reject new task
	if p.closing.Load() {
		return
	}
	p.pending.Inc()
	p.tasks <- task
}

// Shutdown shutdowns all worker goroutines gracefully
func (p *pool) Shutdown() {
	if p.closing.Swap(true) {
		// already closing
		return
	}
	// do shutdown logic
	p.stop <- struct{}{}
	<-p.stop

	// process pending tasks, before shutdown
	completed := make(chan struct{})
	go func() {
		for taskFn := range p.tasks {
			taskFn()
		}
		completed <- struct{}{}
	}()
	// close tasks chan
	close(p.tasks)
	// wait pending task process complete
	<-completed
}

// completeTask completes the task
func (p *pool) completeTask() {
	p.completed.Inc()
	p.worker.Inc()
}

// dispatch dispatches the task to free worker
func (p *pool) dispatch() {
	for {
		select {
		case task := <-p.tasks:
			// get free worker
			worker := <-p.workers
			p.worker.Dec()
			worker.execute(task)
			// dec pending after execute
			p.pending.Dec()
		case <-p.stop:
			// shutdown all worker
			for i := 0; i < cap(p.workers); i++ {
				worker := <-p.workers
				worker.shutdown()
			}
			p.stop <- struct{}{}
			return
		}
	}
}

// worker represents the worker that executes the task
type worker struct {
	pool *pool

	tasks chan Task
	stop  chan struct{}
}

// newWorker creates the worker and start goroutine
func newWorker(pool *pool) *worker {
	w := &worker{
		pool:  pool,
		tasks: make(chan Task),
		stop:  make(chan struct{}),
	}
	go w.process()
	return w
}

// execute submits the task to queue
func (w *worker) execute(task Task) {
	w.tasks <- task
}

// process process task from queue
func (w *worker) process() {
	for {
		w.pool.workers <- w
		select {
		case taskFn := <-w.tasks:
			taskFn()
			w.pool.completeTask()
		case <-w.stop:
			w.stop <- struct{}{}
			log.Info("worker exist...", logger.String("pool", w.pool.name))
			return
		}
	}
}

// shutdown shutdowns the worker goroutine
func (w *worker) shutdown() {
	w.stop <- struct{}{}
	// waiting for goroutine exit
	<-w.stop
}
