// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package annotation

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/apache/arrow-go/v18/arrow/array"
	larray "github.com/lindb/arrow/pkg/arrow/array"
	"github.com/lindb/client_go/api"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/collections"
	tpl "github.com/lindb/lindb/pkg/template"
)

const (
	metricName      = "name"
	metricTimestamp = "timestamp"
	metricTags      = "tags"
	metricFields    = "fields"
)

type MetricMapper struct {
	props *collections.Properties

	nameTpl      *template.Template
	name         string
	nameCol      int
	timestampCol int
	tagCols      []int
	fieldCols    []int

	initialized bool

	logger logger.Logger
}

func NewMetricMapper(props *collections.Properties) Mapper {
	return &MetricMapper{
		props:        props,
		nameCol:      -1,
		timestampCol: -1,
		logger:       logger.GetLogger("CEP", "MetricMapper"),
	}
}

func (m *MetricMapper) Map(event models.Event) models.Event {
	record, ok := event.(arrow.RecordBatch)
	if !ok {
		return nil
	}
	if !m.initialized {
		m.initialize(record)
	}
	fmt.Println(record)

	m.logger.Info("metric mapper...", logger.Any("page", record))

	var points []*api.Point
	for i := 0; i < int(record.NumRows()); i++ {
		point := api.NewPoint(m.getMetricName(record, i)). // metric name
									SetTimestamp(m.getTimestamp(record, i)) // timestamp
		m.buildTags(record, i, point)   // tags
		m.buildFields(record, i, point) // fields

		if m.nameTpl != nil {
			// TODO:: remove tag if be used in metric name template
			var buf bytes.Buffer
			if err := m.nameTpl.Execute(&buf, point.Tags()); err != nil {
				m.logger.Warn("failed to execute metric name template",
					logger.Any("template", m.name), logger.Error(err))
				continue
			}
			point.SetMetricName(buf.String())
		}

		points = append(points, point)

		if m.logger.Enabled(logger.InfoLevel) {
			m.logger.Info("mapped point", logger.Any("name", point.MetricName()),
				logger.Any("timestamp", point.Timestamp()),
				logger.Any("tags", point.Tags()),
				logger.Any("fields", len(point.Fields())))
		}
	}

	return points
}

func (m *MetricMapper) initialize(input arrow.RecordBatch) {
	layout := make(map[string]arrow.Field)
	refs := make(map[string]int)
	for i, col := range input.Schema().Fields() {
		refs[col.Name] = i
		layout[col.Name] = col
	}
	m.name = m.props.GetStringDefault(metricName, "name")

	// support template metric name using golang template
	tpl, err := tpl.Parse("metricName", m.name)
	if err != nil {
		m.logger.Warn("failed to parse metric name template, treat as static name",
			logger.Any("name", m.name), logger.Error(err))
	}
	if tpl != nil {
		m.nameTpl = tpl
	} else {
		m.nameCol = refs[m.name]
	}

	m.timestampCol = refs[m.props.GetStringDefault(metricTimestamp, "timestamp")]
	tags, _ := m.props.GetStringSlice(metricTags)
	m.tagCols = lo.Map(tags, func(tag string, _ int) int {
		return refs[tag]
	})
	fields, _ := m.props.GetStringSlice(metricFields)
	m.fieldCols = lo.Map(fields, func(field string, _ int) int {
		return refs[field]
	})
	m.initialized = true

	m.logger.Info("initialized metric mapper",
		logger.Any("nameCol", m.nameCol),
		logger.Any("timestampCol", m.timestampCol),
		logger.Any("tagCols", m.tagCols),
		logger.Any("fieldCols", m.fieldCols))
}

func (m *MetricMapper) getMetricName(record arrow.RecordBatch, row int) string {
	if m.nameCol != -1 {
		val := record.Column(m.nameCol).(*array.String).Value(row)
		if val != "" {
			return val
		}
	}
	return m.name
}

func (m *MetricMapper) getTimestamp(record arrow.RecordBatch, row int) time.Time {
	if m.timestampCol != -1 {
		val := record.Column(m.timestampCol).(*array.Timestamp).Value(row).ToTime(arrow.Millisecond)
		if !val.IsZero() {
			return val
		}
	}
	// TODO: fast time?
	return time.Now()
}

func (m *MetricMapper) buildTags(record arrow.RecordBatch, row int, point *api.Point) {
	for _, col := range m.tagCols {
		// TODO: add check -1?
		// if col == nil {
		// 	continue
		// }
		column := record.Column(col)
		switch c := column.(type) {
		case *array.String:
			val := c.Value(row)
			if val != "" {
				point.AddTag(record.Schema().Fields()[col].Name, val)
			}
		case *array.Map:
			keys := c.Keys().(*array.String)
			items := c.Items().(*array.String)
			offsets := c.Offsets()
			start, end := offsets[row], offsets[row+1]
			for i := int(start); i < int(end); i++ {
				point.AddTag(keys.Value(i), items.Value(i))
			}
		}
	}
}

func (m *MetricMapper) buildFields(record arrow.RecordBatch, row int, point *api.Point) {
	for _, col := range m.fieldCols {
		// if col == nil {
		// 	continue
		// }
		column := record.Column(col)
		field := record.Schema().Field(col)
		aggType, ok := field.Metadata.GetValue("agg")
		if !ok {
			m.logger.Warn("field column missing agg type metadata, skip",
				logger.Any("column", field.Name))
			continue
		}
		switch aggType {
		case "sum":
			if c, ok := column.(*array.Float64); ok {
				val := c.Value(row)
				point.AddField(api.NewSum(field.Name, val))
			}
		case "first":
			if c, ok := column.(*array.Float64); ok {
				val := c.Value(row)
				point.AddField(api.NewFirst(field.Name, val))
			}
		case "last":
			if c, ok := column.(*array.Float64); ok {
				val := c.Value(row)
				point.AddField(api.NewLast(field.Name, val))
			}
		case "min":
			if c, ok := column.(*array.Float64); ok {
				val := c.Value(row)
				point.AddField(api.NewMin(field.Name, val))
			}
		case "max":
			if c, ok := column.(*array.Float64); ok {
				val := c.Value(row)
				point.AddField(api.NewMax(field.Name, val))
			}
		case "exemplar":
			if c, ok := column.(*larray.Exemplar); ok {
				traceID, spanID, duration := c.Value(row)

				point.AddField(api.NewExemplar(
					field.Name,
					string(traceID),
					string(spanID),
					duration,
				))
			}
		default:
			m.logger.Warn("unsupported agg type for field column, skip",
				logger.Any("column", field.Name), logger.Any("aggType", aggType))
		}

		fmt.Println("field column:", record.Schema().Fields()[col].Name, column.DataType())
	}
}
