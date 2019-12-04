package field

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Fields(t *testing.T) {
	var fs Fields
	fs = append(fs,
		Field{Name: []byte("a"), Type: SumField, Value: float64(0)},
		Field{Name: []byte("c"), Type: HistogramField, Value: float64(0)},
		Field{Name: []byte("b"), Type: SummaryField, Value: float64(0)})
	sort.Sort(fs)

	fs = fs.Insert(Field{Name: []byte("b"), Type: MaxField, Value: float64(0)})
	assert.Equal(t, MaxField, fs[1].Type)

	fs = fs.Insert(Field{Name: []byte("d"), Type: MinField, Value: float64(0)})
	assert.Equal(t, MinField, fs[3].Type)
}
