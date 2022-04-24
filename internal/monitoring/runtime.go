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

package monitoring

import (
	"runtime"

	"github.com/lindb/lindb/internal/linmetric"
	"github.com/lindb/lindb/metrics"
)

type runtimeObserver struct {
	statistics *metrics.RuntimeStatistics
	msLast     *runtime.MemStats
}

func newRuntimeObserver(r *linmetric.Registry) *runtimeObserver {
	observer := &runtimeObserver{
		statistics: metrics.NewRuntimeStatistics(r),
		msLast:     &runtime.MemStats{},
	}
	runtime.ReadMemStats(observer.msLast)
	return observer
}

func (rb *runtimeObserver) Observe() {
	rb.statistics.Routines.Update(float64(runtime.NumGoroutine()))
	n, _ := runtime.ThreadCreateProfile(nil)
	rb.statistics.Threads.Update(float64(n))

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)

	rb.statistics.Alloc.Update(float64(ms.Alloc))
	rb.statistics.TotalAlloc.Add(float64(ms.TotalAlloc - rb.msLast.TotalAlloc))
	rb.statistics.Sys.Update(float64(ms.Sys))
	rb.statistics.LooksUp.Add(float64(ms.Lookups - rb.msLast.Lookups))
	rb.statistics.Mallocs.Add(float64(ms.Mallocs - rb.msLast.Mallocs))
	rb.statistics.Frees.Add(float64(ms.Frees - rb.msLast.Frees))

	rb.statistics.HeapAlloc.Update(float64(ms.HeapAlloc))
	rb.statistics.HeapSys.Update(float64(ms.HeapSys))
	rb.statistics.HeadIdle.Update(float64(ms.HeapIdle))
	rb.statistics.HeapInuse.Update(float64(ms.HeapInuse))
	rb.statistics.HeapReleased.Update(float64(ms.HeapReleased))
	rb.statistics.HeapObjects.Update(float64(ms.HeapObjects))
	rb.statistics.StackInUse.Update(float64(ms.StackInuse))
	rb.statistics.StackSys.Update(float64(ms.StackSys))
	rb.statistics.MSpanInuse.Update(float64(ms.MSpanInuse))
	rb.statistics.MSpanSys.Update(float64(ms.MSpanSys))
	rb.statistics.MCacheInuse.Update(float64(ms.MCacheInuse))
	rb.statistics.MCacheSys.Update(float64(ms.MCacheSys))
	rb.statistics.BuckHashSys.Update(float64(ms.BuckHashSys))
	rb.statistics.GCSys.Update(float64(ms.GCSys))
	rb.statistics.OtherSys.Update(float64(ms.OtherSys))
	rb.statistics.NextGC.Update(float64(ms.NextGC))
	rb.statistics.LastGC.Update(float64(ms.LastGC))
	rb.statistics.GCCPUFraction.Update(ms.GCCPUFraction)

	*rb.msLast = ms
}
