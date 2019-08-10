package task

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const kindDummy Kind = "you-guess"

type dummyParams struct{}

func (p dummyParams) Bytes() []byte { return []byte("{}") }

type dummyProcessor struct{ callcnt int32 }

func (p *dummyProcessor) Kind() Kind                  { return kindDummy }
func (p *dummyProcessor) RetryCount() int             { return 0 }
func (p *dummyProcessor) RetryBackOff() time.Duration { return 0 }
func (p *dummyProcessor) Concurrency() int            { return 1 }
func (p *dummyProcessor) Process(ctx context.Context, task Task) error {
	atomic.AddInt32(&p.callcnt, 1)
	return nil
}
func (p *dummyProcessor) CallCount() int { return int(atomic.LoadInt32(&p.callcnt)) }

func TestTask(t *testing.T) {
	assert.Equal(t, "StateCreated", StateCreated.String())
	assert.Equal(t, "StateRunning", StateRunning.String())
	assert.Equal(t, "StateDoneErr", StateDoneErr.String())
	assert.Equal(t, "StateDoneOK", StateDoneOK.String())
}
