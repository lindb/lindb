package monitoring

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func Test_NewRuntimeCollector(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runtimeCollector := NewRunTimeCollector(
		ctx,
		time.Millisecond*100,
		nil)
	// manually trigger gc to cover ReportCounter
	runtime.GC()

	go runtimeCollector.Run()

	time.Sleep(time.Second * 2)
}
