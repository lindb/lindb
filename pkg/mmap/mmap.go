package mmap

import (
	"syscall"
	"os"
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
	data, err := syscall.Mmap(int(f.Fd()), 0, int(size), syscall.PROT_READ, syscall.MAP_SHARED)
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
	return syscall.Munmap(data)
}
