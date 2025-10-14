package plan

import (
	"fmt"
	"testing"
)

func Test_CleanNameHint(t *testing.T) {
	fmt.Println(cleanNameHint("name_1"))
	fmt.Println(cleanNameHint("name"))
}
