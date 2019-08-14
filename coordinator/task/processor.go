package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/damnever/goctl/queue"
	"github.com/damnever/goctl/retry"
	"github.com/damnever/goctl/semaphore"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./processor.go -destination=./processor_mock.go -package=task

// Processor is responsible for process actual tasks.
// The caller must ensure Process func is goroutine safe if it changes the shared state.
type Processor interface {
	Kind() Kind
	RetryCount() int
	RetryBackOff() time.Duration
	Concurrency() int
	Process(ctx context.Context, task Task) error
}

type taskEvent struct {
	key  string
	task Task
	rev  int64
}

type taskProcessor struct {
	ctx    context.Context
	cancel context.CancelFunc

	repo      state.Repository
	taskq     *queue.Queue
	retrier   retry.Retrier
	sem       *semaphore.Semaphore
	wg        sync.WaitGroup
	processor Processor
}

func newTaskProcessor(ctx context.Context, proc Processor, repo state.Repository) *taskProcessor {
	concurrency := proc.Concurrency()
	if concurrency <= 0 {
		concurrency = 1
	}

	cctx, cancel := context.WithCancel(ctx)
	p := &taskProcessor{
		ctx:    cctx,
		cancel: cancel,

		repo:      repo,
		taskq:     queue.NewQueue(),
		retrier:   retry.New(retry.ConstantBackoffs(proc.RetryCount(), proc.RetryBackOff())),
		sem:       semaphore.NewSemaphore(concurrency),
		processor: proc,
	}
	go p.run()
	return p
}

func (p *taskProcessor) Submit(task taskEvent) error {
	if p.processor.Kind() != task.task.Kind {
		return fmt.Errorf("task name mismatch: %s", task.task.Name)
	}
	select {
	case <-p.ctx.Done(): // High priority..
		return p.ctx.Err()
	default:
		p.taskq.Put(task)
	}
	return nil
}

func (p *taskProcessor) run() {
	for {
		item, err := p.taskq.Get(p.ctx)
		if err != nil {
			return
		}
		if err := p.sem.Acquire(p.ctx); err != nil {
			return
		}
		p.wg.Add(1)
		go p.process(item.(taskEvent))
	}
}

func (p *taskProcessor) process(evt taskEvent) {
	log := logger.GetLogger("coordinator/task/processor")
	defer func() {
		// wait task process done
		p.wg.Done()
		_ = p.sem.Release()
		if e := recover(); e != nil {
			log.Error("process task", logger.Error(fmt.Errorf("panic: %v", e)),
				logger.String("name", evt.key))
		}
	}()

	// Update Result
	err := p.retrier.Run(p.ctx, func() (st retry.State, err error) {
		err = p.processor.Process(p.ctx, evt.task)
		// TODO(damnever): stop if error is fatal
		return
	})
	task := evt.task
	if err != nil {
		log.Error("process task", logger.String("name", evt.key), logger.Error(err))
		task.State = StateDoneErr
		task.ErrMsg = err.Error()
	} else {
		task.State = StateDoneOK
	}

	// save task status
	txn := p.repo.NewTransaction()
	txn.ModRevisionCmp(evt.key, "=", evt.rev)
	txn.Put(evt.key, encoding.JSONMarshal(&task))
	if err := p.repo.Commit(p.ctx, txn); err != nil {
		log.Error("update task status", logger.String("name", evt.key), logger.Error(err))
	}
}

func (p *taskProcessor) Stop() {
	p.cancel()
	p.wg.Wait()
}
