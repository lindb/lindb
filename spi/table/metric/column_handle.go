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

package metric

import (
	"fmt"

	"github.com/lindb/common/models"

	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/sql/tree"
)

type Column interface {
	isExemplar() bool
	createStream(numOfFamilies, numOfPoints int)
	downsampling(familyLoaders []*loader)
	aggregate(aggregator []Result)
	reset()
}

type aggregator[V float64 | *models.Exemplar] struct {
	index int

	start int64
	step  int64

	aggFn  aggregateFunc[V]
	source *rollup[V]
}

func (a *aggregator[V]) aggregate(dst Result) {
	// aggregate down sampled data
	src := a.source.getTimeSeries()
	for index, timestamp := range src.timestamps {
		value := src.values[index]
		offset := 0
		if a.step > 0 {
			offset = int((timestamp - a.start) / a.step)
		}
		if offset < 0 {
			// TODO: add log
			panic(fmt.Sprintf("warn offset < 0, offset=%v\n", offset))
		}
		switch d := dst.(type) {
		case *result[V]:
			if d.array.HasValue(offset) {
				d.array.SetValue(offset, a.aggFn(d.array.GetValue(offset), value))
			} else {
				d.array.SetValue(offset, value)
			}
		}
	}
}

type column[V float64 | *models.Exemplar] struct {
	offset  int
	table   *TableScan
	field   field.Meta
	handles []*ColumnHandle

	streamAgg aggregateFunc[V]
	streams   *Streams[V] // streams of this column, for data family loader

	rollups []*rollup[V]
	aggs    []*aggregator[V]

	numOfPoints int
}

func newColumn[V float64 | *models.Exemplar](offset int,
	table *TableScan,
	field field.Meta,
	handles []*ColumnHandle,
	streamAgg aggregateFunc[V],
	downsmapling getAggregateFunc[V],
) *column[V] {
	c := &column[V]{
		offset:    offset,
		table:     table,
		field:     field,
		handles:   handles,
		streamAgg: streamAgg,
	}
	c.initialize(downsmapling)
	return c
}

func (c *column[V]) isExemplar() bool {
	return c.field.Type.IsExemplar()
}

func (c *column[V]) createStream(numOfFamilies, numOfPoints int) {
	c.streams = newStreams[V](numOfFamilies, numOfPoints)
}

func (c *column[V]) initialize(fn getAggregateFunc[V]) {
	rollupMap := make(map[tree.FuncName]*rollup[V])
	timeRange := c.table.timeRange
	interval := c.table.interval
	step := interval.Int64()
	for i, handle := range c.handles {
		rollup, ok := rollupMap[handle.Downsampling]
		if !ok {
			rollup = newRollup(timeRange.NumOfPoints(interval), step, fn(handle.Downsampling))
			// maybe duplicate downsampling func
			rollupMap[handle.Downsampling] = rollup
			c.rollups = append(c.rollups, rollup)
		}
		c.aggs = append(c.aggs, &aggregator[V]{
			index:  c.offset + i,
			start:  timeRange.Start,
			step:   step,
			source: rollup,
			aggFn:  fn(handle.Aggregation),
		})
	}

	if c.table.isTimestampSelected {
		c.numOfPoints = timeRange.NumOfPoints(interval)
	} else {
		c.numOfPoints = 1 // timestamp not in select item list
	}
}

func (c *column[V]) load(familyIndex int, slotRange timeutil.SlotRange, getter func(slot uint16) (V, bool)) {
	columnStream := c.streams.GetStreamByIndex(familyIndex, func(numOfPoints int) Stream[V] {
		return NewStream[V](numOfPoints)
	})

	fn := c.streamAgg

	for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
		value, ok := getter(movingSourceSlot)
		if !ok {
			// no data, goto next loop
			continue
		}
		columnStream.SetAtStep(int(movingSourceSlot), value, fn)
	}
}

func (c *column[V]) downsampling(familyLoaders []*loader) {
	for familyIndex, loader := range familyLoaders {
		familyTime := loader.familyTime
		slotRange := loader.timeRange
		interval := loader.interval.Int64()
		for movingSourceSlot := slotRange.Start; movingSourceSlot <= slotRange.End; movingSourceSlot++ {
			timestamp := familyTime + int64(movingSourceSlot)*interval
			value := c.streams.GetStreamByIndex(familyIndex, func(numOfPoints int) Stream[V] {
				return NewStream[V](numOfPoints)
			}).GetAtStep(int(movingSourceSlot))

			// rollup
			for _, r := range c.rollups {
				r.doRollup(timestamp, value)
			}
		}
	}
}

func (c *column[V]) aggregate(aggregator []Result) {
	for _, agg := range c.aggs {
		if aggregator[agg.index] == nil {
			aggregator[agg.index] = NewResult[V](c.numOfPoints)
		}
		agg.aggregate(aggregator[agg.index])
	}
}

// reset resets the column rollup states for reusing context.
func (c *column[V]) reset() {
	// reset streams
	c.streams.reset()

	// reset rollups
	for _, r := range c.rollups {
		r.reset()
	}
}
