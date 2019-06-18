package kv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	var lock = NewLock("t.lock")
	var err = lock.Lock()
	assert.Nil(t, err, "lock error")

	err = lock.Lock()
	assert.NotNil(t, err, "cannot lock again for locked file")

	err = lock.Unlock()
	assert.Nil(t, err, "unlock error")

	lock = NewLock("t.lock")
	err = lock.Lock()
	assert.Nil(t, err, "lock error")

	lock.Unlock()

	fileInfo, _ := os.Stat("t.lock")
	assert.Nil(t, fileInfo, "lock file exist")

	lock = NewLock("/tmp/not_dir/t.lock")
	err = lock.Lock()
	assert.NotNil(t, err, "cannot lock not exist file")
}
