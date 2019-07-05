package lockers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_FileLocker(t *testing.T) {
	fl1 := NewFileLocker("/lindb/storage/1")
	assert.True(t, fl1.TryLock())
	assert.False(t, fl1.TryLock())

	fl2 := NewFileLocker("/lindb/storage/2")
	assert.True(t, fl2.TryLock())

	fl3 := NewFileLocker("/lindb/storage/2")
	assert.False(t, fl3.TryLock())

	fl2.Unlock()
	assert.True(t, fl2.TryLock())

}
