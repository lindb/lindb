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

package metricsdata

import (
	"errors"
	"sort"

	"github.com/lindb/common/models"
	"github.com/lindb/common/pkg/logger"
	"github.com/lindb/roaring"
	"github.com/samber/lo"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/collections"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

var mergerLogger = logger.GetLogger("MetricData", "Merger")

var MetricDataMerger kv.MergerType = "MetricDataMerger"

// init registers metric data merger create function
func init() {
	kv.RegisterMerger(MetricDataMerger, NewMerger)
}

type mergerContext struct {
	// per metric level context
	scanners                 []*dataScanner
	seriesIDs                *roaring.Bitmap // target series ids
	allFields                field.Metas     // all target fields(normal+exemplar)
	targetRange, sourceRange timeutil.SlotRange
	fieldReader              FieldReader

	// global context
	ratio             uint16
	baseSlot          uint16
	targetNumOfPoints int

	// merger context per series
	decoder   *encoding.TSDDecoder
	values    *collections.Array[float64]
	exemplars *collections.Array[*models.Exemplar]
}

func (mc *mergerContext) reset() {
	mc.scanners = mc.scanners[:0]
	mc.seriesIDs.Clear()
	mc.allFields = mc.allFields[:0]
}

// merger implements kv.Merger for merging series data for each metric
type merger struct {
	dataFlusher  Flusher
	seriesMerger SeriesMerger
	rollup       kv.Rollup
	option       kv.FamilyOption

	context *mergerContext
}

// NewMerger creates a metric data merger
func NewMerger() kv.Merger {
	return &merger{}
}

// Init initializes metric data merger, if rollup context exist do rollup job, else do compact job
func (m *merger) Init(flusher kv.Flusher, params map[string]any) error {
	dataFlusher, err := NewFlusher(flusher)
	if err != nil {
		return err
	}
	option, ok := params[kv.FamilyOptionContext]
	if !ok {
		return errors.New("missing family option context")
	}
	familyOption := option.(kv.FamilyOption)
	m.option = familyOption

	if rollupCtx, ok := params[kv.RollupContext]; ok {
		m.rollup = rollupCtx.(kv.Rollup)
	}
	m.dataFlusher = dataFlusher
	m.seriesMerger = newSeriesMerger(dataFlusher)

	m.context = &mergerContext{
		seriesIDs: roaring.New(),
		decoder:   encoding.GetTSDDecoder(),
	}

	// check if rollup job
	if m.rollup != nil {
		m.context.ratio = m.rollup.IntervalRatio()
		m.context.baseSlot = m.rollup.BaseSlot() // different family, need calc based on base slot
		m.context.targetNumOfPoints = m.rollup.NumOfPoints()
	} else {
		m.context.ratio = 1
		m.context.targetNumOfPoints = m.option.NumOfPoints
	}

	m.context.values = collections.NewArray[float64](m.context.targetNumOfPoints)
	m.context.exemplars = collections.NewArray[*models.Exemplar](m.context.targetNumOfPoints)
	return nil
}

// Merge merges the multi metric data into one target metric data for same metric id
func (m *merger) Merge(key uint32, metricBlocks [][]byte) error {
	defer m.context.reset()

	// 1. prepare readers and metric level data(field/time slot/series ids)
	err := m.prepare(metricBlocks)
	if err != nil {
		return err
	}
	// All input blocks were corrupted and skipped — nothing to merge for this key.
	// Return nil so the compact job continues with the next key instead of failing.
	if len(m.context.scanners) == 0 {
		return nil
	}
	// 2. Prepare metric
	m.dataFlusher.PrepareMetric(key, m.context.allFields)
	// 3. merge series data by roaring container
	highKeys := m.context.seriesIDs.GetHighKeys()

	for idx, highKey := range highKeys {
		container := m.context.seriesIDs.GetContainerAtIndex(idx)
		it := container.PeekableIterator()
		for it.HasNext() {
			lowSeriesID := it.Next()
			if err := m.seriesMerger.merge(m.context, highKey, lowSeriesID); err != nil {
				return err
			}
			// flush series id
			if err := m.dataFlusher.FlushSeries(encoding.ValueWithHighLowBits(uint32(highKey)<<16, lowSeriesID)); err != nil {
				return err
			}
		}
	}
	// flush metric data
	if err := m.dataFlusher.CommitMetric(m.context.targetRange); err != nil {
		return err
	}
	return nil
}

func (m *merger) prepare(metricBlocks [][]byte) error {
	for _, metricBlock := range metricBlocks {
		reader, err := NewReader("merge_operation", metricBlock)
		if err != nil {
			// Corrupted block — log and skip; the rest of the blocks are still mergeable.
			// Returning an error here would abort the entire compact job permanently.
			mergerLogger.Warn("skipping corrupted metric block during merge", logger.Error(err))
			continue
		}
		m.context.seriesIDs.Or(reader.GetSeriesIDs())
		// get target slot range(start/end)
		timeRange := reader.GetTimeRange()
		if len(m.context.allFields) == 0 {
			m.context.sourceRange.Start = timeRange.Start
			m.context.sourceRange.End = timeRange.End
		} else {
			m.context.sourceRange = m.context.sourceRange.Union(timeRange)
		}
		// merge target fields under metric level
		m.context.allFields = append(m.context.allFields, reader.GetFields()...)
		// create data scanner
		scanner, err := newDataScanner(reader)
		if err != nil {
			return err
		}
		m.context.scanners = append(m.context.scanners, scanner)
	}

	// deduplicate target fields
	lo.UniqBy(m.context.allFields, func(f field.Meta) field.ID { return f.ID })

	// sort by field id
	sort.Slice(m.context.allFields, func(i, j int) bool { return m.context.allFields[i].ID < m.context.allFields[j].ID })

	// check if rollup job
	if m.rollup != nil {
		// calc target time slot range and interval ratio
		m.context.targetRange.Start = m.rollup.CalcSlot(m.rollup.GetTimestamp(m.context.sourceRange.Start))
		m.context.targetRange.End = m.rollup.CalcSlot(m.rollup.GetTimestamp(m.context.sourceRange.End))
	} else {
		m.context.targetRange.Start = m.context.sourceRange.Start
		m.context.targetRange.End = m.context.sourceRange.End
	}
	return nil
}
