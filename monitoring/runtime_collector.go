package monitoring

import (
	"context"
	"io"
	"os"
	"runtime"
	"time"

	"github.com/uber-go/tally"
)

// RunTimeCollector collects the go-runtime information
type RunTimeCollector struct {
	ctx         context.Context
	lastGCCount uint32
	interval    time.Duration
	closer      io.Closer
	scope       tally.Scope
}

// NewRunTimeCollector returns a new collector to retrieve runtime metrics
func NewRunTimeCollector(
	ctx context.Context,
	interval time.Duration,
	tags map[string]string,
) *RunTimeCollector {
	host, _ := os.Hostname()
	if tags == nil {
		tags = make(map[string]string)
	}
	tags["host"] = host

	scope, closer := tally.NewRootScope(tally.ScopeOptions{
		Tags:   tags,
		Prefix: "runtime",
	}, interval)

	return &RunTimeCollector{
		ctx:      ctx,
		scope:    scope,
		closer:   closer,
		interval: interval,
	}
}

func (c *RunTimeCollector) Run() {
	defer func() {
		_ = c.closer.Close()
	}()

	c.collect()
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.collect()
		case <-c.ctx.Done():
			return
		}
	}
}

func (c *RunTimeCollector) collect() {
	var stat runtime.MemStats
	runtime.ReadMemStats(&stat)

	// General statistics.
	generalScope := c.scope.SubScope("mem")
	generalScope.Gauge("alloc").Update(float64(stat.Alloc))
	generalScope.Gauge("total_alloc").Update(float64(stat.TotalAlloc))
	generalScope.Gauge("sys").Update(float64(stat.Sys))
	generalScope.Gauge("lookups").Update(float64(stat.Lookups))
	generalScope.Gauge("mallocs").Update(float64(stat.Mallocs))
	generalScope.Gauge("frees").Update(float64(stat.Frees))

	// Heap memory statistics.
	heapScope := c.scope.SubScope("heap")
	heapScope.Gauge("alloc").Update(float64(stat.HeapAlloc))
	heapScope.Gauge("sys").Update(float64(stat.HeapSys))
	heapScope.Gauge("idle").Update(float64(stat.HeapIdle))
	heapScope.Gauge("inuse").Update(float64(stat.HeapInuse))
	heapScope.Gauge("released").Update(float64(stat.HeapReleased))
	heapScope.Gauge("objects").Update(float64(stat.HeapObjects))

	// Stack memory statistics.
	stackScope := c.scope.SubScope("stack")
	stackScope.Gauge("inuse").Update(float64(stat.StackInuse))
	stackScope.Gauge("sys").Update(float64(stat.StackSys))

	// Off-heap memory statistics.
	offHeapScope := c.scope.SubScope("offheap")
	offHeapScope.Gauge("mspan_inuse").Update(float64(stat.MSpanInuse))
	offHeapScope.Gauge("mspan_sys").Update(float64(stat.MSpanSys))
	offHeapScope.Gauge("mcache_inuse").Update(float64(stat.MCacheInuse))
	offHeapScope.Gauge("mcache_sys").Update(float64(stat.MCacheSys))
	offHeapScope.Gauge("buck_hash_sys").Update(float64(stat.BuckHashSys))

	// Garbage collector statistics.
	gcScope := c.scope.SubScope("gc")
	gcScope.Gauge("next").Update(float64(stat.NextGC))
	gcScope.Gauge("last_interval").Update(float64(stat.LastGC))
	gcScope.Gauge("pause_total_ns").Update(float64(stat.PauseTotalNs))
	gcScope.Counter("count").Inc(int64(500 + stat.NumGC - c.lastGCCount))
	c.lastGCCount = stat.NumGC

	// Goroutines statistics
	c.scope.Gauge("num_goroutines").Update(float64(runtime.NumGoroutine()))
}
