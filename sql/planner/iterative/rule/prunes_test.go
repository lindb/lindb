package rule

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/sql/planner/plan"
)

func TestA(t *testing.T) {
	fmt.Println(reflect.TypeFor[*plan.OutputNode]())
	fmt.Println(reflect.TypeFor[plan.OutputNode]())
	fmt.Println(reflect.TypeOf((&plan.OutputNode{})))
	assert.True(t, reflect.TypeFor[*plan.OutputNode]() == reflect.TypeOf(&plan.OutputNode{}))
}
