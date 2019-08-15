// +build windows

package lockers

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	procLockFileEx = syscall.NewLazyDLL("kernel32.dll").NewProc("LockFileEx")

	errLocked = errors.New("The process cannot access the file because another process has locked a portion of the file.")
)

const (
	// https://msdn.microsoft.com/en-us/library/windows/desktop/aa365203(v=vs.85).aspx
	LOCKFILE_EXCLUSIVE_LOCK   = 2
	LOCKFILE_FAIL_IMMEDIATELY = 1

	// see https://msdn.microsoft.com/en-us/library/windows/desktop/ms681382(v=vs.85).aspx
	errLockViolation syscall.Errno = 0x21
)

func lockFile(path string, flag int) (*os.File, error) {
	f, err := open(path, flag)
	if err != nil {
		return nil, err
	}
	if err := lockFileShort(syscall.Handle(f.Fd()), 0); err != nil {
		f.Close()
		return nil, err
	}
	return f, nil
}

func unlockFile(f *os.File) error {
	return unlockFileEx(syscall.Handle(f.Fd()), 1, 0, &syscall.Overlapped{})
}

func open(path string, flag int) (*os.File, error) {
	if path == "" {
		return nil, fmt.Errorf("cannot open empty filename")
	}
	var access uint32
	switch flag {
	case syscall.O_RDONLY:
		access = syscall.GENERIC_READ
	case syscall.O_WRONLY:
		access = syscall.GENERIC_WRITE
	case syscall.O_RDWR:
		access = syscall.GENERIC_READ | syscall.GENERIC_WRITE
	case syscall.O_WRONLY | syscall.O_CREAT:
		access = syscall.GENERIC_ALL
	default:
		panic(fmt.Errorf("flag %v is not supported", flag))
	}
	fd, err := syscall.CreateFile(&(syscall.StringToUTF16(path)[0]),
		access,
		syscall.FILE_SHARE_READ|syscall.FILE_SHARE_WRITE|syscall.FILE_SHARE_DELETE,
		nil,
		syscall.OPEN_ALWAYS,
		syscall.FILE_ATTRIBUTE_NORMAL,
		0)
	if err != nil {
		return nil, err
	}
	return os.NewFile(uintptr(fd), path), nil
}

func lockFileShort(fd syscall.Handle, flags uint32) error {
	flag := uint32(LOCKFILE_EXCLUSIVE_LOCK)
	flag |= flags
	if fd == syscall.InvalidHandle {
		return nil
	}
	err := lockFileEx(fd, flag, 1, 0, &syscall.Overlapped{})
	if err == nil {
		return nil
	} else if err.Error() == errLocked.Error() {
		return errLocked
	} else if err != errLockViolation {
		return err
	}
	return nil
}

func lockFileEx(h syscall.Handle, flags, locklow, lockhigh uint32, ol *syscall.Overlapped) (err error) {
	var reserved uint32
	r1, _, e1 := syscall.Syscall6(
		procLockFileEx.Addr(),
		6,
		uintptr(h),
		uintptr(flags),
		uintptr(reserved),
		uintptr(locklow),
		uintptr(lockhigh),
		uintptr(unsafe.Pointer(ol)),
	)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}

func unlockFileEx(h syscall.Handle, locklow, lockhigh uint32, ol *syscall.Overlapped) (err error) {
	var reserved uint32 = 0
	r1, _, e1 := syscall.Syscall6(
		procLockFileEx.Addr(),
		5,
		uintptr(h),
		uintptr(reserved),
		uintptr(locklow),
		uintptr(lockhigh),
		uintptr(unsafe.Pointer(ol)),
		0,
	)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return err
}
