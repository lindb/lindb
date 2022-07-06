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

package metrics

import "github.com/lindb/lindb/internal/linmetric"

// RuntimeStatistics represents golang runtime statistics.
type RuntimeStatistics struct {
	Routines     *linmetric.BoundGauge   // the number of goroutines
	Threads      *linmetric.BoundGauge   // the number of records in the thread creation profile
	Alloc        *linmetric.BoundGauge   // bytes of allocated heap objects
	TotalAlloc   *linmetric.BoundCounter // cumulative bytes allocated for heap objects
	Sys          *linmetric.BoundGauge   // the total bytes of memory obtained from the OS
	LooksUp      *linmetric.BoundCounter // the number of pointer lookups performed by the runtime
	Mallocs      *linmetric.BoundCounter // the cumulative count of heap objects allocated
	Frees        *linmetric.BoundCounter // the cumulative count of heap objects freed
	HeapAlloc    *linmetric.BoundGauge   // bytes of allocated heap objects
	HeapSys      *linmetric.BoundGauge   // bytes of heap memory obtained from the OS
	HeadIdle     *linmetric.BoundGauge   // bytes in idle (unused) spans
	HeapInuse    *linmetric.BoundGauge   // bytes in in-use spans
	HeapReleased *linmetric.BoundGauge   // bytes of physical memory returned to the OS
	HeapObjects  *linmetric.BoundGauge   // the number of allocated heap objects
	StackInUse   *linmetric.BoundGauge   // bytes in stack spans
	StackSys     *linmetric.BoundGauge   // bytes of stack memory obtained from the OS
	MSpanInuse   *linmetric.BoundGauge   // bytes of allocated mspan structures
	MSpanSys     *linmetric.BoundGauge   // bytes of memory obtained from the OS for mspan
	MCacheInuse  *linmetric.BoundGauge   // bytes of allocated mcache structures
	MCacheSys    *linmetric.BoundGauge   // bytes of memory obtained from the OS for mcache structures
	BuckHashSys  *linmetric.BoundGauge   // bytes of memory in profiling bucket hash tables
	GCSys        *linmetric.BoundGauge   // bytes of memory in garbage collection metadata
	OtherSys     *linmetric.BoundGauge   // bytes of memory in miscellaneous off-heap
	NextGC       *linmetric.BoundGauge   // the target heap size of the next GC cycle
	LastGC       *linmetric.BoundGauge   // the time the last garbage collection finished
	// the fraction of this program's available  CPU time used by the GC since the program started.
	GCCPUFraction *linmetric.BoundGauge
}

// NewRuntimeStatistics creates a golang runtime statistics.
func NewRuntimeStatistics(registry *linmetric.Registry) *RuntimeStatistics {
	runtimeScope := registry.NewScope("lindb.runtime")
	memoryScope := runtimeScope.Scope("mem")
	return &RuntimeStatistics{
		Routines:      runtimeScope.NewGauge("go_goroutines"),
		Threads:       runtimeScope.NewGauge("go_threads"),
		Alloc:         memoryScope.NewGauge("alloc"),
		TotalAlloc:    memoryScope.NewCounter("total_alloc"),
		Sys:           memoryScope.NewGauge("sys"),
		LooksUp:       memoryScope.NewCounter("lookups"),
		Mallocs:       memoryScope.NewCounter("mallocs"),
		Frees:         memoryScope.NewCounter("frees"),
		HeapAlloc:     memoryScope.NewGauge("heap_alloc"),
		HeapSys:       memoryScope.NewGauge("heap_sys"),
		HeadIdle:      memoryScope.NewGauge("heap_idle"),
		HeapInuse:     memoryScope.NewGauge("heap_inuse"),
		HeapReleased:  memoryScope.NewGauge("heap_released"),
		HeapObjects:   memoryScope.NewGauge("heap_objects"),
		StackInUse:    memoryScope.NewGauge("stack_inuse"),
		StackSys:      memoryScope.NewGauge("stack_sys"),
		MSpanInuse:    memoryScope.NewGauge("mspan_inuse"),
		MSpanSys:      memoryScope.NewGauge("mspan_sys"),
		MCacheInuse:   memoryScope.NewGauge("mcache_inuse"),
		MCacheSys:     memoryScope.NewGauge("mcache_sys"),
		BuckHashSys:   memoryScope.NewGauge("buck_hash_sys"),
		GCSys:         memoryScope.NewGauge("gc_sys"),
		OtherSys:      memoryScope.NewGauge("other_sys"),
		NextGC:        memoryScope.NewGauge("next_gc"),
		LastGC:        memoryScope.NewGauge("last_gc"),
		GCCPUFraction: memoryScope.NewGauge("gc_cpu_fraction"),
	}
}
