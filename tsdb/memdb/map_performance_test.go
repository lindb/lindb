package memdb

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/eleme/lindb/pkg/lockers"
)

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
	index := key & 255
	sm.maps[index].mu.RLock()
	v, ok := sm.maps[index].m[index]
	sm.maps[index].mu.RUnlock()
	return v, ok
}

func (sm *shardingRwLockedMap) Set(key int, value int) {
	index := key & 255
	sm.maps[index].mu.Lock()
	sm.maps[index].m[key] = value
	sm.maps[index].mu.Unlock()
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
