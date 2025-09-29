package trace

import (
	"fmt"
	"testing"
)

func TestTraceMergeOperator(t *testing.T) {
	merger := &TraceMergeOperator{}
	dest, _ := merger.FullMerge([]byte("key"), []byte{1, 2, 3}, [][]byte{{4, 5, 6}, {8, 9}})
	fmt.Println(dest)
}
