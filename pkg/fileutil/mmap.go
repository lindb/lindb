package fileutil

import (
	"os"

	"go.uber.org/zap"

	"github.com/eleme/lindb/pkg/logger"
)

var log = logger.GetLogger("pkg/mmap")

const (
	read = 1 << iota
	write
)

// Map memory-maps a file.
func Map(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error("mmap map close file error", zap.Error(err))
		}
	}()

	fs, err := f.Stat()
	if err != nil {
		return nil, err
	}
	size := fs.Size()
	if size == 0 {
		return nil, nil
	}

	// map file
	data, err := mmap(int(f.Fd()), 0, int(size), read)

	if err != nil {
		return nil, err
	}
	return data, nil
}

// RWMap maps a file for read and write with give size.
// New file is created is not existed.
func RWMap(path string, size int) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Error("mmap rwmap close file error", zap.Error(err))
		}
	}()

	fstat, err := f.Stat()

	if err != nil {
		return nil, err
	}

	if fstat.Size() < int64(size) {
		if err := f.Truncate(int64(size)); err != nil {
			return nil, err
		}
	}

	// map file
	data, err := mmap(int(f.Fd()), 0, size, read|write)

	if err != nil {
		return nil, err
	}
	return data, nil
}

// Unmap closes the memory-map.
func Unmap(data []byte) error {
	if data == nil {
		return nil
	}
	return munmap(data)
}

func Sync(data []byte) error {
	return msync(data)
}
