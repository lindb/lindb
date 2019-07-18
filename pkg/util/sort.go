package util

import "sort"

//SortUint32 sorts a slice of uint32 in increasing order.
func SortUint32(array []uint32) []uint32 {
	sort.Slice(array, func(i, j int) bool {
		return array[i] < array[j]
	})
	return array
}
