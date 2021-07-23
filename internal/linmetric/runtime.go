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

package linmetric

import "runtime"

type runtimeObserver struct {
	goRoutinesGauge        *BoundGauge
	threadsGauge           *BoundGauge
	allocBytesGauge        *BoundGauge
	allocBytesTotalCounter *BoundDeltaCounter
	sysBytesGauge          *BoundGauge
	looksUpCounter         *BoundDeltaCounter
	mallocsTotalCounter    *BoundDeltaCounter
	freesTotalCounter      *BoundDeltaCounter
	heapAllocBytesGauge    *BoundGauge
	heapSysBytesGauge      *BoundGauge
	headIdleGague          *BoundGauge
	heapInuseGauge         *BoundGauge
	heapReleasedGauge      *BoundGauge
	heapObjectsGauge       *BoundGauge
	stackInUseGague        *BoundGauge
	stackSysBytes          *BoundGauge
	mspanInuseBytes        *BoundGauge
	mspanSysuseBytes       *BoundGauge
	mcacheInuseBytes       *BoundGauge
	mcacheSysBytes         *BoundGauge
	buckHashSysBytes       *BoundGauge
	gcSysBytes             *BoundGauge
	otherSysBytes          *BoundGauge
	nextGCBytes            *BoundGauge
	lastGCTimeSeconds      *BoundGauge
	gcCPUFraction          *BoundGauge

	msLast *runtime.MemStats
}

func newRuntimeObserver() *runtimeObserver {
	runtimeScope := NewScope("lindb.runtime")
	memoryScope := runtimeScope.Scope("mem")

	observer := &runtimeObserver{
		goRoutinesGauge:        runtimeScope.NewGauge("go_goroutines"),
		threadsGauge:           runtimeScope.NewGauge("go_threads"),
		allocBytesGauge:        memoryScope.NewGauge("alloc_bytes"),
		allocBytesTotalCounter: memoryScope.NewDeltaCounter("alloc_bytes_total"),
		sysBytesGauge:          memoryScope.NewGauge("sys_bytes"),
		looksUpCounter:         memoryScope.NewDeltaCounter("lookups_total"),
		mallocsTotalCounter:    memoryScope.NewDeltaCounter("mallocs_total"),
		freesTotalCounter:      memoryScope.NewDeltaCounter("frees_total"),
		heapAllocBytesGauge:    memoryScope.NewGauge("heap_alloc_bytes"),
		heapSysBytesGauge:      memoryScope.NewGauge("heap_sys_bytes"),
		headIdleGague:          memoryScope.NewGauge("heap_idle_bytes"),
		heapInuseGauge:         memoryScope.NewGauge("heap_inuse_bytes"),
		heapReleasedGauge:      memoryScope.NewGauge("heap_released_bytes"),
		heapObjectsGauge:       memoryScope.NewGauge("heap_objects"),
		stackInUseGague:        memoryScope.NewGauge("stack_inuse_bytes"),
		stackSysBytes:          memoryScope.NewGauge("stack_sys_bytes"),
		mspanInuseBytes:        memoryScope.NewGauge("mspan_inuse_bytes"),
		mspanSysuseBytes:       memoryScope.NewGauge("mspan_sys_bytes"),
		mcacheInuseBytes:       memoryScope.NewGauge("mcache_inuse_bytes"),
		mcacheSysBytes:         memoryScope.NewGauge("mcache_sys_bytes"),
		buckHashSysBytes:       memoryScope.NewGauge("buck_hash_sys_bytes"),
		gcSysBytes:             memoryScope.NewGauge("gc_sys_bytes"),
		otherSysBytes:          memoryScope.NewGauge("other_sys_bytes"),
		nextGCBytes:            memoryScope.NewGauge("next_gc_bytes"),
		lastGCTimeSeconds:      memoryScope.NewGauge("last_gc_time_seconds"),
		gcCPUFraction:          memoryScope.NewGauge("gc_cpu_fraction"),
		msLast:                 &runtime.MemStats{},
	}
	runtime.ReadMemStats(observer.msLast)
	return observer
}

func (rb *runtimeObserver) Observe() {
	rb.goRoutinesGauge.Update(float64(runtime.NumGoroutine()))
	n, _ := runtime.ThreadCreateProfile(nil)
	rb.threadsGauge.Update(float64(n))

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	rb.allocBytesGauge.Update(float64(ms.Alloc))
	rb.allocBytesTotalCounter.Add(float64(ms.TotalAlloc - rb.msLast.TotalAlloc))
	rb.sysBytesGauge.Update(float64(ms.Sys))
	rb.looksUpCounter.Add(float64(ms.Lookups - rb.msLast.Lookups))
	rb.mallocsTotalCounter.Add(float64(ms.Mallocs - rb.msLast.Mallocs))
	rb.freesTotalCounter.Add(float64(ms.Frees - rb.msLast.Frees))

	rb.heapAllocBytesGauge.Update(float64(ms.HeapAlloc))
	rb.heapSysBytesGauge.Update(float64(ms.HeapSys))
	rb.headIdleGague.Update(float64(ms.HeapIdle))
	rb.heapInuseGauge.Update(float64(ms.HeapInuse))
	rb.heapReleasedGauge.Update(float64(ms.HeapReleased))
	rb.heapObjectsGauge.Update(float64(ms.HeapObjects))
	rb.stackInUseGague.Update(float64(ms.StackInuse))
	rb.stackSysBytes.Update(float64(ms.StackSys))
	rb.mspanInuseBytes.Update(float64(ms.MSpanInuse))
	rb.mspanSysuseBytes.Update(float64(ms.MSpanSys))
	rb.mcacheInuseBytes.Update(float64(ms.MCacheInuse))
	rb.mcacheSysBytes.Update(float64(ms.MCacheSys))
	rb.buckHashSysBytes.Update(float64(ms.BuckHashSys))
	rb.gcSysBytes.Update(float64(ms.GCSys))
	rb.otherSysBytes.Update(float64(ms.OtherSys))
	rb.nextGCBytes.Update(float64(ms.NextGC))
	rb.lastGCTimeSeconds.Update(float64(ms.LastGC))
	rb.gcCPUFraction.Update(ms.GCCPUFraction)

	*rb.msLast = ms
}
