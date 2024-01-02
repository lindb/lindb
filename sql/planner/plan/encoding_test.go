package plan

import (
	"fmt"
	"reflect"
	"testing"
)

func TestEncoding(t *testing.T) {
	fmt.Println(reflect.TypeOf(TableScanNode{}).String())
	fmt.Println(reflect.TypeOf(&TableScanNode{}).Elem())
}
