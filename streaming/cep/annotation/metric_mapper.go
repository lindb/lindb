package annotation

import (
	"fmt"

	"github.com/samber/lo"

	"github.com/lindb/lindb/spi/types"
)

type MetricMapper struct {
	name      string
	timestamp string
	tags      []string
	fields    []string

	nameCol      int
	timestampCol int
	tagCols      []int
	fieldCols    []int

	initialized bool
}

func NewMetricMapper() *MetricMapper {
	return &MetricMapper{
		name:      "metric_name",
		timestamp: "timestamp",
		tags:      []string{"tag", "interface"},
		fields:    []string{"qps"},

		nameCol:      -1,
		timestampCol: -1,
		tagCols:      []int{},
		fieldCols:    []int{},
	}
}

func (m *MetricMapper) Map(input *types.Page) (any, error) {
	if !m.initialized {
		m.initialize(input)
	}
	it := input.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		fmt.Printf("metric name: %v\n", row.GetString(m.nameCol))
		fmt.Printf("timestamp: %v\n", row.GetTimestamp(m.timestampCol))
		tags := lo.Map(m.tagCols, func(col int, _ int) string {
			return fmt.Sprintf("%v", row.Get(col))
		})
		fmt.Printf("tags: %v\n", tags)
		fields := lo.Map(m.fieldCols, func(col int, _ int) int64 {
			return int64(*row.GetInt(col))
		})
		fmt.Printf("fields: %v\n", fields)
	}

	return nil, nil
}

func (m *MetricMapper) initialize(input *types.Page) {
	layout := make(map[string]int)
	for i, col := range input.Layout {
		layout[col.Name] = i
	}
	m.nameCol = layout[m.name]
	m.timestampCol = layout[m.timestamp]
	m.tagCols = lo.Map(m.tags, func(tag string, _ int) int {
		return layout[tag]
	})
	m.fieldCols = lo.Map(m.fields, func(field string, _ int) int {
		return layout[field]
	})
	fmt.Println(m.tagCols)

	m.initialized = true
}
