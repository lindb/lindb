// +build amd64

package trie

import (
	"golang.org/x/sys/cpu"
)

var hasBMI2 = cpu.X86.HasBMI2

// go:noescape
func select64(x uint64, k int64) int64
