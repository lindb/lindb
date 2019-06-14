package storage

import (
	"testing"
	"os"
)

func Test_Lock(t *testing.T) {
	var lock = NewLock("t.lock")
	var err = lock.Lock()
	if nil != err {
		t.Error(err)
	}

	err = lock.Lock()
	if nil == err {
		t.Errorf("can't lock locked file again")
	}

	err = lock.Unlock()
	if nil != err {
		t.Error(err)
	}

	lock = NewLock("t.lock")
	err = lock.Lock()
	if nil != err {
		t.Error(err)
	}

	lock.Unlock()

	fileInfo, err := os.Stat("t.lock")

	if nil != fileInfo {
		t.Errorf("lock file exist")
	}

	lock = NewLock("/tmp/not_dir/t.lock")
	err = lock.Lock()
	if nil == err {
		t.Errorf("fail: create file succuess")
	}
}
