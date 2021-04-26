// +build (NOT windows)

package main

import "fmt"

func getSpace(string) (int, int, int, error) {
	return 0, 0, 0, fmt.Errorf("not implemented")
}
