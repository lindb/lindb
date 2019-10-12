package concurrent

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/atomic"
)

func TestPool_Submit(t *testing.T) {
	grNum := runtime.NumGoroutine()
	pool := NewPool("test_pool", 2, 10)
	assert.Equal(t, "test_pool", pool.Name())
	// num. of pool + 1 dispatch
	assert.Equal(t, grNum+2+1, runtime.NumGoroutine())

	var c atomic.Int32
	iterations := 10000

	for i := 0; i < iterations; i++ {
		pool.Execute(func() {
			c.Inc()
		})
	}
	pool.Shutdown()
	pool.Shutdown()
	// reject all task
	for i := 0; i < iterations; i++ {
		pool.Execute(func() {
			c.Inc()
		})
	}
	assert.Equal(t, int32(iterations), c.Load())
	// all goroutines of pool exited
	assert.Equal(t, grNum, runtime.NumGoroutine())
}
