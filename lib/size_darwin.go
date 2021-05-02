package lib

import "fmt"

// GetSpace returns disk size information for the given drive (bytes).
func GetSpace(string) (uint64, uint64, uint64, error) {
	return 0, 0, 0, fmt.Errorf("not implemented")
}
