package utils

import (
	"fmt"
	"testing"

	"github.com/lindb/lindb/spi/types"
)

func Test_Values(t *testing.T) {
	page := types.NewPage()
	v1 := NewValues(page)
	v2 := NewValues(page)
	fmt.Println(v1.GetID())
	fmt.Println(v2.GetID())
}
