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

package metricsdata

import (
	"sort"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

var MetricDataMerger kv.MergerType = "MetricDataMerger"

// init registers metric data merger create function
func init() {
	kv.RegisterMerger(MetricDataMerger, NewMerger)
}

type mergerContext struct {
	scanners     []*dataScanner
	seriesIDs    *roaring.Bitmap // target series ids
	targetFields field.Metas     // target fields

	targetRange, sourceRange timeutil.SlotRange
	ratio                    uint16
}

// merger implements kv.Merger for merging series data for each metric
type merger struct {
	dataFlusher  Flusher
	seriesMerger SeriesMerger
	rollup       kv.Rollup
}

// NewMerger creates a metric data merger
func NewMerger(flusher kv.Flusher) (kv.Merger, error) {
	dataFlusher, err := NewFlusher(flusher)
	if err != nil {
		return nil, err
	}
	return &merger{
		dataFlusher:  dataFlusher,
		seriesMerger: newSeriesMerger(dataFlusher),
	}, nil
}

// Init initializes metric data merger, if rollup context exist do rollup job, else do compact job
func (m *merger) Init(params map[string]interface{}) {
	rollupCtx, ok := params[kv.RollupContext]
	if ok {
		m.rollup = rollupCtx.(kv.Rollup)
	}
}

// Merge merges the multi metric data into one target metric data for same metric id
func (m *merger) Merge(key uint32, metricBlocks [][]byte) error {
	blockCount := len(metricBlocks)
	// 1. prepare readers and metric level data(field/time slot/series ids)
	mergeCtx, err := m.prepare(metricBlocks)
	if err != nil {
		return err
	}
	// 2. Prepare metric
	m.dataFlusher.PrepareMetric(key, mergeCtx.targetFields)
	// 3. merge series data by roaring container
	highKeys := mergeCtx.seriesIDs.GetHighKeys()
	decodeStreams := make([]*encoding.TSDDecoder, blockCount) // make decodeStreams for reuse
	defer func() {
		for _, stream := range decodeStreams {
			encoding.ReleaseTSDDecoder(stream)
		}
	}()
	encodeStream := encoding.TSDEncodeFunc(mergeCtx.targetRange.Start)
	fieldReaders := make([]FieldReader, blockCount)
	for idx, highKey := range highKeys {
		container := mergeCtx.seriesIDs.GetContainerAtIndex(idx)
		it := container.PeekableIterator()
		for it.HasNext() {
			lowSeriesID := it.Next()
			// maybe series id not exist in some value block
			for blockIdx, scanner := range mergeCtx.scanners {
				seriesEntry := scanner.scan(highKey, lowSeriesID)
				if len(seriesEntry) == 0 {
					continue
				}
				timeRange := scanner.slotRange()
				if fieldReaders[blockIdx] == nil {
					fieldReaders[blockIdx] = newFieldReader(scanner.fieldIndexes(), seriesEntry, timeRange)
				} else {
					fieldReaders[blockIdx].Reset(seriesEntry, timeRange)
				}
			}
			if err := m.seriesMerger.merge(mergeCtx, decodeStreams, encodeStream, fieldReaders); err != nil {
				return err
			}
			// flush series id
			if err := m.dataFlusher.FlushSeries(encoding.ValueWithHighLowBits(uint32(highKey)<<16, lowSeriesID)); err != nil {
				return err
			}
		}
	}
	// flush metric data
	if err := m.dataFlusher.CommitMetric(mergeCtx.targetRange); err != nil {
		return err
	}
	return nil
}

func (m *merger) prepare(metricBlocks [][]byte) (*mergerContext, error) {
	ctx := &mergerContext{
		scanners:     make([]*dataScanner, len(metricBlocks)),
		seriesIDs:    roaring.New(),
		targetFields: field.Metas{},
	}

	for idx, metricBlock := range metricBlocks {
		reader, err := NewReader("merge_operation", metricBlock)
		if err != nil {
			return nil, err
		}
		ctx.seriesIDs.Or(reader.GetSeriesIDs())
		// get target slot range(start/end)
		timeRange := reader.GetTimeRange()
		if len(ctx.targetFields) == 0 {
			ctx.sourceRange.Start = timeRange.Start
			ctx.sourceRange.End = timeRange.End
		} else {
			if ctx.sourceRange.Start > timeRange.Start {
				ctx.sourceRange.Start = timeRange.Start
			}
			if ctx.sourceRange.End < timeRange.End {
				ctx.sourceRange.End = timeRange.End
			}
		}
		// merge target fields under metric level
		for _, f := range reader.GetFields() {
			_, ok := ctx.targetFields.GetFromID(f.ID)
			if !ok {
				ctx.targetFields = ctx.targetFields.Insert(f)
			}
		}
		// create data scanner
		if ctx.scanners[idx], err = newDataScanner(reader); err != nil {
			return nil, err
		}
	}
	// sort by field id
	sort.Slice(ctx.targetFields, func(i, j int) bool { return ctx.targetFields[i].ID < ctx.targetFields[j].ID })
	// check if rollup job

	if m.rollup != nil {
		// calc target time slot range and interval ratio
		ctx.targetRange.Start = m.rollup.CalcSlot(m.rollup.GetTimestamp(ctx.sourceRange.Start))
		ctx.targetRange.End = m.rollup.CalcSlot(m.rollup.GetTimestamp(ctx.sourceRange.End))
		ctx.ratio = m.rollup.IntervalRatio()
	} else {
		ctx.targetRange.Start = ctx.sourceRange.Start
		ctx.targetRange.End = ctx.sourceRange.End
		ctx.ratio = 1
	}
	return ctx, nil
}
