package ingest

import (
	"fmt"
	"testing"
)

func TestLog(t *testing.T) {
	for i := range 10 {
		fmt.Printf("%v\n", i)
	}
}
