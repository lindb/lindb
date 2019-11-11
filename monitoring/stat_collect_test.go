package monitoring

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/lindb/lindb/models"
)

type mockReport struct {
	wg     *sync.WaitGroup
	cancel context.CancelFunc
}

func (r *mockReport) Report(stat interface{}) {
	r.cancel()
	r.wg.Done()
}

func TestNewStatusReporter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	var wg sync.WaitGroup
	wg.Add(1)

	_ = NewStatCollect(ctx, 10*time.Millisecond, "/tmp", &mockReport{wg: &wg, cancel: cancel}, models.ActiveNode{})
	wg.Wait()
}
