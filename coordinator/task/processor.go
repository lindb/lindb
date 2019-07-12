package task

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/damnever/goctl/queue"
	"github.com/damnever/goctl/retry"
	"github.com/damnever/goctl/semaphore"
	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
	"github.com/eleme/lindb/pkg/state"
)

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

	cli       state.Repository
	taskq     *queue.Queue
	retrier   retry.Retrier
	sem       *semaphore.Semaphore
	wg        sync.WaitGroup
	processor Processor
}

func newTaskProcessor(ctx context.Context, proc Processor, cli state.Repository) *taskProcessor {
	concurrency := proc.Concurrency()
	if concurrency <= 0 {
		concurrency = 1
	}

	cctx, cancel := context.WithCancel(ctx)
	p := &taskProcessor{
		ctx:    cctx,
		cancel: cancel,

		cli:       cli,
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
	log := logger.GetLogger()
	defer func() {
		p.wg.Done()
		_ = p.sem.Release()
		if e := recover(); e != nil {
			log.Error("process task", zap.Error(fmt.Errorf("panic: %v", e)),
				zap.String("name", evt.key))
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
		log.Error("process task", zap.String("name", evt.key), zap.Error(err))
		task.State = StateDoneErr
		task.ErrMsg = err.Error()
	} else {
		task.State = StateDoneOK
	}
	//TODO ????????modify status
	//resp, err := p.cli.Txn(p.ctx).If(
	//	etcdcliv3.Compare(etcdcliv3.ModRevision(evt.key), "=", evt.rev),
	//).Then(
	//	etcdcliv3.OpPut(evt.key, zerocopy.UnsafeBtoa(task.UnsafeMarshal())),
	//).Commit()
	//if err := state.TxnErr(resp, err); err != nil {
	//	log.Error("update task status", zap.String("name", evt.key), zap.Error(err))
	//}
}

func (p *taskProcessor) Stop() {
	p.cancel()
	p.wg.Wait()
}
