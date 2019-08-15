package fileutil

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLockedFile(t *testing.T) {
	path := filepath.Join(os.TempDir(), "linfl")
	defer os.RemoveAll(path)

	f1, err := LockFile(path, os.O_CREATE|os.O_RDWR, 0644)
	require.Nil(t, err)
	_, err = TryLockFile(path, os.O_CREATE|os.O_RDWR, 0644)
	require.NotNil(t, err)
	require.Nil(t, f1.Close())

	f2, err := LockFile(path, os.O_CREATE|os.O_RDWR, 0644)
	require.Nil(t, err)
	donec := make(chan struct{})
	var f3 *LockedFile
	go func() {
		var err error
		f3, err = LockFile(path, os.O_CREATE|os.O_RDWR, 0644)
		require.Nil(t, err)
		close(donec)
	}()
	select {
	case <-donec:
		t.Fatal("multiple locks on same file")
	case <-time.After(time.Second):
	}
	require.Nil(t, f2.Close())
	select {
	case <-donec:
	case <-time.After(time.Second):
		t.Fatal("cannot acquire lock on file")
	}
	_, err = TryLockFile(path, os.O_CREATE|os.O_RDWR, 0644)
	require.NotNil(t, err)
	require.Nil(t, f3.Close())
}
