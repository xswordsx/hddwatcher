package lib

import (
	"fmt"
	"syscall"
	"unsafe"
)

// GetSpace returns disk size information for the given drive (bytes).
func GetSpace(drive string) (avail int64, total int64, free int64, err error) {
	kernel32, err := syscall.LoadLibrary("Kernel32.dll")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot load kernel32.dll: %w", err)
	}
	defer syscall.FreeLibrary(kernel32)
	GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot find GetDiskFreeSpaceExW: %w", err)
	}

	lpFreeBytesAvailable := int64(0)
	lpTotalNumberOfBytes := int64(0)
	lpTotalNumberOfFreeBytes := int64(0)
	ok, _, msg := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive))),
		uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
		uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

	if ok == 0 {
		return 0, 0, 0, fmt.Errorf("cannot get disk space: %w", msg)
	}

	return lpFreeBytesAvailable, lpTotalNumberOfBytes, lpTotalNumberOfFreeBytes, nil
}
