package memdb

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/eleme/lindb/models"
	"github.com/eleme/lindb/pkg/field"
	"github.com/eleme/lindb/pkg/lockers"
	"github.com/eleme/lindb/pkg/timeutil"
	"github.com/eleme/lindb/tsdb/index"
	"github.com/eleme/lindb/tsdb/metrictbl"

	"github.com/golang/mock/gomock"
)

///////////////////////////////////////////////////
//                mock interface
///////////////////////////////////////////////////

func makeMockIDGenerator(ctrl *gomock.Controller) *index.MockIDGenerator {
	mockGen := index.NewMockIDGenerator(ctrl)
	mockGen.EXPECT().GenTSID(gomock.Any(), gomock.Any()).
		Return(uint32(2222)).AnyTimes()
	mockGen.EXPECT().GenFieldID(gomock.Any(), gomock.Any(), gomock.Any()).
		Return(uint32(1111)).AnyTimes()
	mockGen.EXPECT().GenMetricID(gomock.Any()).
		Return(uint32(3333)).AnyTimes()

	return mockGen
}

func makeMockTableWriter(ctrl *gomock.Controller) *metrictbl.MockTableWriter {
	mockTW := metrictbl.NewMockTableWriter(ctrl)
	mockTW.EXPECT().WriteField(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return().AnyTimes()
	mockTW.EXPECT().WriteTSEntry(gomock.Any()).
		Return().AnyTimes()
	mockTW.EXPECT().WriteMetricBlock(gomock.Any()).
		Return(nil).AnyTimes()
	mockTW.EXPECT().Close().Return(nil).AnyTimes()

	return mockTW
}

func makeMockPoint(ctrl *gomock.Controller) *models.MockPoint {
	mockPoint := models.NewMockPoint(ctrl)

	mockPoint.EXPECT().Name().Return("cpu.load").AnyTimes()
	mockPoint.EXPECT().Tags().Return("idle").AnyTimes()
	mockPoint.EXPECT().Timestamp().Return(timeutil.Now()).AnyTimes()

	fakeFields := make(map[string]models.Field)
	fakeHistogram := models.NewMockField(ctrl)
	fakeHistogram.EXPECT().Type().Return(field.HistogramField).AnyTimes()
	fakeHistogram.EXPECT().IsComplex().Return(true).AnyTimes()
	fakeFields["histogram"] = fakeHistogram

	fakeMax := models.NewMockSimpleField(ctrl)
	fakeMax.EXPECT().ValueType().Return(field.Integer).AnyTimes()
	fakeMax.EXPECT().AggType().Return(field.Sum).AnyTimes()
	fakeMax.EXPECT().Type().Return(field.SumField).AnyTimes()
	fakeMax.EXPECT().Value().Return(1).AnyTimes()
	fakeMax.EXPECT().IsComplex().Return(false).AnyTimes()

	fakeFields["max"] = fakeMax

	mockPoint.EXPECT().Fields().Return(fakeFields).AnyTimes()
	return mockPoint
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
