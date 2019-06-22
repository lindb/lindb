package lockers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileLock(t *testing.T) {
	var lock = NewFileLock("t.lock")
	var err = lock.Lock()
	assert.Nil(t, err, "lock error")

	err = lock.Lock()
	assert.NotNil(t, err, "cannot lock again for locked file")

	err = lock.Unlock()
	assert.Nil(t, err, "unlock error")

	lock = NewFileLock("t.lock")
	err = lock.Lock()
	assert.Nil(t, err, "lock error")

	lock.Unlock()

	fileInfo, _ := os.Stat("t.lock")
	assert.Nil(t, fileInfo, "lock file exist")

	lock = NewFileLock("/tmp/not_dir/t.lock")
	err = lock.Lock()
	assert.NotNil(t, err, "cannot lock not exist file")
}
