package lockers

import (
	"sync"
)

// filesMap is a simple mutex locked map containing files in use.
var filesMap = &struct {
	sync.Mutex
	files map[string]struct{}
}{
	files: make(map[string]struct{}),
}

// FileLocker provides the ability of restricting access to a specified file for a single process.
// thread-safe, not process-safe.
type FileLocker interface {
	// TryLock will try to lock the file and return whether it succeed or not without blocking.
	TryLock() bool
	// Unlock unlocks the file, this operation is reentrant。
	Unlock()
}

// fileLocker implements FileLocker
type fileLocker struct {
	fileName string
}

// NewFileLocker returns a new FileLocker.
func NewFileLocker(fileName string) FileLocker {
	return &fileLocker{fileName: fileName}
}

// TryLock try locking file, return false if locked.
func (fl *fileLocker) TryLock() bool {
	filesMap.Lock()
	defer filesMap.Unlock()

	_, ok := filesMap.files[fl.fileName]
	if ok {
		return false
	}
	filesMap.files[fl.fileName] = struct{}{}
	return true
}

// Unlock unlocks file lock.
func (fl *fileLocker) Unlock() {
	filesMap.Lock()
	defer filesMap.Unlock()

	delete(filesMap.files, fl.fileName)
}
