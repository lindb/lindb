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
	"html/template"
	"time"

	"github.com/lindb/client_go/api"
	"github.com/lindb/common/pkg/logger"
	"github.com/samber/lo"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/collections"
	tpl "github.com/lindb/lindb/pkg/template"
	"github.com/lindb/lindb/spi/types"
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
	nameCol      *types.ColumnMetadata
	timestampCol *types.ColumnMetadata
	tagCols      []*types.ColumnMetadata
	fieldCols    []*types.ColumnMetadata

	initialized bool

	logger logger.Logger
}

func NewMetricMapper(props *collections.Properties) Mapper {
	return &MetricMapper{
		props:  props,
		logger: logger.GetLogger("CEP", "MetricMapper"),
	}
}

func (m *MetricMapper) Map(event models.Event) models.Event {
	page, ok := event.(*types.Page)
	if !ok {
		return nil
	}
	if !m.initialized {
		m.initialize(page)
	}
	var points []*api.Point
	it := page.Iterator()
	for row := it.Begin(); row != it.End(); row = it.Next() {
		point := api.NewPoint(m.getMetricName(row)). // metric name
								SetTimestamp(m.getTimestamp(row)) // timestamp
		m.buildTags(row, point)   // tags
		m.buildFields(row, point) // fields

		// TODO: add exemplars support
		point.AddField(api.NewExemplar("exemplar", "traceid", "spanid", 10))

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

func (m *MetricMapper) initialize(input *types.Page) {
	layout := make(map[string]*types.ColumnMetadata)
	for i, col := range input.Layout {
		col.Ref = i
		layout[col.Name] = &col
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
		m.nameCol = layout[m.name]
	}

	m.timestampCol = layout[m.props.GetStringDefault(metricTimestamp, "timestamp")]
	tags, _ := m.props.GetStringSlice(metricTags)
	m.tagCols = lo.Map(tags, func(tag string, _ int) *types.ColumnMetadata {
		return layout[tag]
	})
	fields, _ := m.props.GetStringSlice(metricFields)
	m.fieldCols = lo.Map(fields, func(field string, _ int) *types.ColumnMetadata {
		return layout[field]
	})
	m.initialized = true

	m.logger.Info("initialized metric mapper",
		logger.Any("nameCol", m.nameCol),
		logger.Any("timestampCol", m.timestampCol),
		logger.Any("tagCols", m.tagCols),
		logger.Any("fieldCols", m.fieldCols))
}

func (m *MetricMapper) getMetricName(row types.Row) string {
	if m.nameCol != nil {
		val := row.GetString(m.nameCol.Ref)
		if val != "" {
			return val
		}
	}
	return m.name
}

func (m *MetricMapper) getTimestamp(row types.Row) time.Time {
	if m.timestampCol != nil {
		val := row.GetTimestamp(m.timestampCol.Ref)
		if !val.IsZero() {
			return val
		}
	}
	return time.Now()
}

func (m *MetricMapper) buildTags(row types.Row, point *api.Point) {
	for _, col := range m.tagCols {
		if col == nil {
			continue
		}
		switch col.DataType {
		case types.DTString:
			val := row.GetString(col.Ref)
			if val != "" {
				point.AddTag(col.Name, val)
			}
		case types.DTMap:
			val := row.GetMap(col.Ref)
			for k, v := range val {
				point.AddTag(k, v)
			}
		}
	}
}

func (m *MetricMapper) buildFields(row types.Row, point *api.Point) {
	for _, col := range m.fieldCols {
		if col == nil {
			continue
		}
		val := float64(row.GetInt(col.Ref))
		switch col.AggType {
		case types.ATSum:
			point.AddField(api.NewSum(col.Name, val))
		case types.ATFirst:
			point.AddField(api.NewFirst(col.Name, val))
		case types.ATLast:
			point.AddField(api.NewLast(col.Name, val))
		case types.ATMin:
			point.AddField(api.NewMin(col.Name, val))
		case types.ATMax:
			point.AddField(api.NewMax(col.Name, val))
		default:
			panic("implement me set field based on agg type")
		}
	}
}
