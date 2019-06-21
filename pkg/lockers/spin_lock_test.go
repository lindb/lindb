package lockers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSpinLock_TryLock(t *testing.T) {
	var sl SpinLock
	assert.True(t, sl.TryLock())
	assert.False(t, sl.TryLock())
}

func TestSpinLock_Lock_Unlock(t *testing.T) {
	var sl SpinLock
	sl.Lock()
	assert.False(t, sl.TryLock())

	sl.Unlock()
	assert.True(t, sl.TryLock())
	sl.Unlock()

	sl.Lock()
	go func() {
		time.Sleep(time.Millisecond)
		sl.Unlock()
	}()
	sl.Lock()
	assert.False(t, sl.TryLock())
}
