package storage

import (
	"os"
	"testing"
)

func Test_Lock(t *testing.T) {
	var lock = NewLock("t.lock")
	var err = lock.Lock()
	if nil != err {
		t.Error(err)
		return
	}

	err = lock.Lock()
	if nil == err {
		t.Errorf("can't lock locked file again")
		return
	}

	err = lock.Unlock()
	if nil != err {
		t.Error(err)
		return
	}

	lock = NewLock("t.lock")
	err = lock.Lock()
	if nil != err {
		t.Error(err)
		return
	}

	lock.Unlock()

	fileInfo, _ := os.Stat("t.lock")

	if nil != fileInfo {
		t.Errorf("lock file exist")
		return
	}

	lock = NewLock("/tmp/not_dir/t.lock")
	err = lock.Lock()
	if nil == err {
		t.Errorf("fail: create file succuess")
		return
	}
}
