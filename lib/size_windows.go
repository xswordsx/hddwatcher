package lib

import (
	"fmt"
	"syscall"
	"unsafe"
)

// GetSpace returns disk size information for the given drive (bytes).
func GetSpace(drive string) (avail uint64, total uint64, free uint64, err error) {
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot load kernel32.dll: %w", err)
	}
	defer syscall.FreeLibrary(kernel32) // nolint: errcheck
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot find GetDiskFreeSpaceExW: %w", err)
	}
	strPtr, err := syscall.UTF16PtrFromString(drive)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot create UTF16 pointer: %w", err)
	}
	lpFreeBytesAvailable := uint64(0)
	lpTotalNumberOfBytes := uint64(0)
	lpTotalNumberOfFreeBytes := uint64(0)
	ok, _, msg := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(strPtr)),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

	if ok == 0 {
		return 0, 0, 0, fmt.Errorf("cannot get disk space: %w", msg)
	}

	return lpFreeBytesAvailable, lpTotalNumberOfBytes, lpTotalNumberOfFreeBytes, nil
}
