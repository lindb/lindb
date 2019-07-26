package fileutil

import (
	"errors"
	"os"
	"reflect"
	"sync"
	"unsafe"

	"golang.org/x/sys/windows"
)

type mapHandle struct {
	file     windows.Handle
	view     windows.Handle
	writable bool
}

var handleMap = make(map[uintptr]*mapHandle)

var lock4map sync.Mutex

func header(bytes []byte) *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
}

func addressAndSize(bytes []byte) (uintptr, uintptr) {
	h := header(bytes)
	return h.Data, uintptr(h.Len)
}

// todo: @TianliangXia test on windows
func mmap(fd int, offset int64, size int, mode int) ([]byte, error) {
	prot := windows.PAGE_READONLY
	access := windows.FILE_MAP_READ
	writable := false
	if mode&write == 1 {
		prot = windows.PAGE_READWRITE
		access = windows.FILE_MAP_WRITE
		writable = true
	}

	// The maximum size is the area of the file, starting from 0,
	// that we wish to allow to be mappable. It is the sum of
	// the length the user requested, plus the offset where that length
	// is starting from. This does not map the data into memory.
	maxSizeHigh := uint32((offset + int64(size)) >> 32)
	maxSizeLow := uint32((offset + int64(size)) & 0xFFFFFFFF)
	// TODO: Do we need to set some security attributes? It might help portability.
	h, errno := windows.CreateFileMapping(windows.Handle(fd), nil, uint32(prot), maxSizeHigh, maxSizeLow, nil)
	if h == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", errno)
	}

	// Actually map a view of the data into memory. The view's size
	// is the length the user requested.
	fileOffsetHigh := uint32(offset >> 32)
	fileOffsetLow := uint32(offset & 0xFFFFFFFF)
	addr, errno := windows.MapViewOfFile(h, uint32(access), fileOffsetHigh, fileOffsetLow, uintptr(size))
	if addr == 0 {
		return nil, os.NewSyscallError("MapViewOfFile", errno)
	}

	lock4map.Lock()
	defer lock4map.Unlock()
	handleMap[addr] = &mapHandle{
		file:     windows.Handle(fd),
		view:     h,
		writable: writable,
	}

	bytes := make([]byte, 0)

	hd := header(bytes)
	hd.Data = addr
	hd.Len = size
	hd.Cap = hd.Len

	return bytes, nil
}

func munmap(bytes []byte) error {
	hd := header(bytes)
	addr := hd.Data
	// Lock the UnmapViewOfFile along with the handleMap deletion.
	// As soon as we unmap the view, the OS is free to give the
	// same addr to another new map. We don't want another goroutine
	// to insert and remove the same addr into handleMap while
	// we're trying to remove our old addr/handle pair.
	lock4map.Lock()
	defer lock4map.Unlock()
	err := windows.UnmapViewOfFile(addr)
	if err != nil {
		return err
	}

	handle, ok := handleMap[addr]
	if !ok {
		// should be impossible; we would've errored above
		return errors.New("unknown base address")
	}
	delete(handleMap, addr)

	e := windows.CloseHandle(windows.Handle(handle.view))
	return os.NewSyscallError("CloseHandle", e)
}

func msync(bytes []byte) error {
	addr, size := addressAndSize(bytes)
	errno := windows.FlushViewOfFile(addr, size)
	if errno != nil {
		return os.NewSyscallError("FlushViewOfFile", errno)
	}

	lock4map.Lock()
	defer lock4map.Unlock()
	handle, ok := handleMap[addr]
	if !ok {
		// should be impossible; we would've errored above
		return errors.New("unknown base address")
	}

	if handle.writable {
		if err := windows.FlushFileBuffers(handle.file); err != nil {
			return os.NewSyscallError("FlushFileBuffers", err)
		}
	}

	return nil
}
