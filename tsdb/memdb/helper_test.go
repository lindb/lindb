package memdb

import (
	"math/rand"
	"runtime"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/pkg/lockers"
	"github.com/lindb/lindb/tsdb/metadb"
	"github.com/lindb/lindb/tsdb/tblstore"
)

///////////////////////////////////////////////////
//                mock interface
///////////////////////////////////////////////////

func makeMockIDGenerator(ctrl *gomock.Controller) *metadb.MockIDGenerator {
	mockGen := metadb.NewMockIDGenerator(ctrl)
	mockGen.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(uint16(1111), nil).AnyTimes()
	mockGen.EXPECT().GenMetricID(gomock.Any()).
		Return(uint32(3333)).AnyTimes()
	mockGen.EXPECT().GenTagKeyID(gomock.Any(), gomock.Any()).Return(uint32(3333)).AnyTimes()
	return mockGen
}

func makeMockDataFlusher(ctrl *gomock.Controller) *tblstore.MockMetricsDataFlusher {
	mockTF := tblstore.NewMockMetricsDataFlusher(ctrl)
	mockTF.EXPECT().FlushFieldMetas(gomock.Any()).Return().AnyTimes()
	mockTF.EXPECT().FlushField(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return().AnyTimes()
	mockTF.EXPECT().FlushSeries(gomock.Any()).
		Return().AnyTimes()
	mockTF.EXPECT().FlushMetric(gomock.Any()).
		Return(nil).AnyTimes()
	mockTF.EXPECT().Commit().Return(nil).AnyTimes()

	return mockTF
}

///////////////////////////////////////////////////
//                benchmark test
///////////////////////////////////////////////////

var _testSyncMap = sync.Map{}

type rwLockedMap struct {
	mu sync.RWMutex
	m  map[int]int
}

type spinLockedMap struct {
	sl lockers.SpinLock
	m  map[int]int
}

func (slm *spinLockedMap) Get(key int) (int, bool) {
	slm.sl.Lock()
	v, ok := slm.m[key]
	slm.sl.Unlock()
	return v, ok

}

func (m *rwLockedMap) Get(key int) (int, bool) {
	m.mu.RLock()
	v, ok := m.m[key]
	m.mu.RUnlock()
	return v, ok
}

type shardingRwLockedMap struct {
	maps [256]rwLockedMap
}

func (sm *shardingRwLockedMap) Get(key int) (int, bool) {
	idx := key & 255
	sm.maps[idx].mu.RLock()
	v, ok := sm.maps[idx].m[idx]
	sm.maps[idx].mu.RUnlock()
	return v, ok
}

func (sm *shardingRwLockedMap) Set(key int, value int) {
	idx := key & 255
	sm.maps[idx].mu.Lock()
	sm.maps[idx].m[key] = value
	sm.maps[idx].mu.Unlock()
}

func Benchmark_syncMap(b *testing.B) {
	for i := 0; i < 10000; i++ {
		_testSyncMap.Store(i, i)
	}

	wg := sync.WaitGroup{}
	for g := 0; g < runtime.NumCPU()*100; g++ {
		wg.Add(1)
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for i := 0; i < b.N; i++ {
				_testSyncMap.Load(r.Intn(10000))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func Benchmark_rwLockedMap(b *testing.B) {
	rwmap := rwLockedMap{m: make(map[int]int)}
	for i := 0; i < 10000; i++ {
		rwmap.m[i] = i
	}
	wg := sync.WaitGroup{}
	for g := 0; g < runtime.NumCPU()*100; g++ {
		wg.Add(1)
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for i := 0; i < b.N; i++ {
				rwmap.Get(r.Intn(10000))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func Benchmark_shardingRwLockedMap(b *testing.B) {
	srwmap := shardingRwLockedMap{}
	for i := 0; i < 256; i++ {
		srwmap.maps[i] = rwLockedMap{m: make(map[int]int)}
	}
	for i := 0; i < 1000; i++ {
		srwmap.Set(i, i)
	}
	wg := sync.WaitGroup{}
	for g := 0; g < runtime.NumCPU()*100; g++ {
		wg.Add(1)
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for i := 0; i < b.N; i++ {
				srwmap.Get(r.Intn(10000))
			}
			wg.Done()
		}()
	}
	wg.Wait()

}

func Benchmark_spinLockMap(b *testing.B) {
	slMap := spinLockedMap{m: make(map[int]int)}
	for i := 0; i < 10000; i++ {
		slMap.m[i] = i
	}

	wg := sync.WaitGroup{}
	for g := 0; g < runtime.NumCPU()*100; g++ {
		wg.Add(1)
		go func() {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))

			for i := 0; i < b.N; i++ {
				slMap.Get(r.Intn(10000))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func Benchmark_100000_get_map(b *testing.B) {
	m := make(map[int]struct{})
	for i := 0; i < 100000; i++ {
		m[i] = struct{}{}
	}

	for x := 0; x < b.N; x++ {
		_ = m[1]
	}
}

func Benchmark_100000_get_slice(b *testing.B) {
	var m []int
	for i := 0; i < 100000; i++ {
		m = append(m, i)
	}
	for x := 0; x < b.N; x++ {
		idx := sort.Search(len(m), func(z int) bool {
			return m[z] >= 1
		})
		_ = m[idx]
	}
}

func Benchmark_100000_map_iterate(b *testing.B) {
	m := make(map[int]struct{})
	for i := 0; i < 100000; i++ {
		m[i] = struct{}{}
	}

	for x := 0; x < b.N; x++ {
		for k, v := range m {
			_, _ = k, v
		}
	}
}

func Benchmark_100000_slice_iterate(b *testing.B) {
	var m []int
	for i := 0; i < 100000; i++ {
		m = append(m, i)
	}
	for x := 0; x < b.N; x++ {
		for k, v := range m {
			_, _ = k, v
		}
	}
}
