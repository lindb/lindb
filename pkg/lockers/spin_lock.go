package lockers

import (
	"runtime"
	"sync/atomic"
)

// SpinLock implements sync/Locker, default 0 indicates an unlocked spinLock.
type SpinLock struct {
	_flag uint32
}

// Lock locks spinLock. If the lock is locked before, the caller will be blocked until unlocked.
func (sl *SpinLock) Lock() {
	for !sl.TryLock() {
		runtime.Gosched() //allow other goroutines to do stuff.
	}
}

// Unlock unlocks spinLock, this operation is reentrant。
func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl._flag, 0)
}

// TryLock will try to lock spinLock and return whether it succeed or not without blocking.
func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl._flag, 0, 1)
}
