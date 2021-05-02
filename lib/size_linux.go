package lib

import (
	"fmt"

	"golang.org/x/sys/unix"
)

// GetSpace returns disk size information for the given drive (bytes).
func GetSpace(drive string) (free uint64, total uint64, avail uint64, err error) {
	var stat unix.Statfs_t
	err = unix.Statfs(drive, &stat)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("cannot stat %q: %w", drive, err)
	}

	// blocks * size per block = available space in bytes
	size := uint64(stat.Bsize)

	total = stat.Blocks * size
	avail = stat.Bavail * size
	free = stat.Bfree * size

	return avail, total, free, nil
}
