package fileutil

import (
	"os"
)

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
	defer f.Close()

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
func RWMap(filePath string, size int) ([]byte, error) {
	f, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	defer f.Close()

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
