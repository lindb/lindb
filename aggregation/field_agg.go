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

package aggregation

import (
	"github.com/lindb/lindb/aggregation/selector"
	"github.com/lindb/lindb/series"
)

//go:generate mockgen -source=./field_agg.go -destination=./field_agg_mock.go -package=aggregation

// FieldAggregator represents a field aggregator, aggregator the field series which with same field id.
type FieldAggregator interface {
	// Aggregate aggregates the field series into current aggregator
	Aggregate(it series.FieldIterator)
	// GetBlock returns series block for saving loaded data
	GetBlock(idx int, fn newBlockFunc) (block series.Block, ok bool)
	// ResultSet returns the result set of field aggregator
	ResultSet() (startTime int64, it series.FieldIterator)

	// reset resets the aggregate context for reusing
	reset()
}

//// aggKey represents aggregate key for supporting multi aggregate function with a same primitive field id
//type aggKey struct {
//	primitiveID field.PrimitiveID
//	aggType     field.AggType
//}

// downSamplingFieldAggregator represents
type downSamplingFieldAggregator struct {
	//segmentStartTime int64
	//start            int

	blockSize int
	blocks    []series.Block
}

// NewDownSamplingFieldAggregator creates a field aggregator for down sampling,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewDownSamplingFieldAggregator(
	aggSpec AggregatorSpec,
	blockSize int,
) FieldAggregator {
	//start, _ := selector.Range()
	agg := &downSamplingFieldAggregator{
		//segmentStartTime: segmentStartTime,
		//start:            start,
		blockSize: blockSize,
		blocks:    make([]series.Block, blockSize),
	}
	//aggregatorMap := make(map[aggKey]PrimitiveAggregator)
	//// if down sampling spec need init all aggregator
	//for funcType := range aggSpec.Functions() {
	//	primitiveFields := aggSpec.GetFieldType().GetPrimitiveFields(funcType)
	//	for _, pField := range primitiveFields {
	//		key := aggKey{
	//			primitiveID: pField.FieldID,
	//			aggType:     pField.AggType,
	//		}
	//		aggregatorMap[key] = NewPrimitiveAggregator(pField.FieldID, selector, pField.AggType.AggFunc())
	//	}
	//}
	//length := len(aggregatorMap)
	//agg.aggregators = make([]PrimitiveAggregator, length)
	//idx := 0
	//for _, pAgg := range aggregatorMap {
	//	agg.aggregators[idx] = pAgg
	//	idx++
	//}
	//// sort field ids
	//sort.Slice(agg.aggregators, func(i, j int) bool {
	//	return agg.aggregators[i].FieldID() < agg.aggregators[j].FieldID()
	//})
	return agg
}

// Aggregate aggregates the field series into current aggregator
func (agg *downSamplingFieldAggregator) Aggregate(_ series.FieldIterator) {
	// do nothing for down sampling
}

// GetBlock returns series block for saving loaded data
func (agg *downSamplingFieldAggregator) GetBlock(idx int, fn newBlockFunc) (block series.Block, ok bool) {
	if idx < 0 || idx >= agg.blockSize {
		return
	}
	block = agg.blocks[idx]
	if block == nil {
		block = fn()
		agg.blocks[idx] = block
	}
	ok = true
	return
}

func (agg *downSamplingFieldAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	//its := make([]series.PrimitiveIterator, len(agg.aggregators))
	//idx := 0
	//for _, it := range agg.aggregators {
	//	its[idx] = it.Iterator()
	//	idx++
	//}
	//return agg.segmentStartTime, newFieldIterator(agg.start, its)
	return
}

// reset resets the aggregate context for reusing
func (agg *downSamplingFieldAggregator) reset() {
	for _, block := range agg.blocks {
		if block != nil {
			block.Clear()
		}
	}
}

// fieldAggregator implements field aggregator interface, aggregator field series based on aggregator spec
type fieldAggregator struct {
	segmentStartTime int64
	start            int

	selector selector.SlotSelector
}

// NewFieldAggregator creates a field aggregator,
// time range 's start and end is index based on segment start time and interval.
// e.g. segment start time = 20190905 10:00:00, start = 10, end = 50, interval = 10 seconds,
// real query time range {20190905 10:01:40 ~ 20190905 10:08:20}
func NewFieldAggregator(segmentStartTime int64, selector selector.SlotSelector) FieldAggregator {
	start, _ := selector.Range()
	agg := &fieldAggregator{
		segmentStartTime: segmentStartTime,
		start:            start,
		selector:         selector,
	}

	return agg
}

// ResultSet returns the result set of field aggregator
func (a *fieldAggregator) ResultSet() (startTime int64, it series.FieldIterator) {
	//its := make([]series.PrimitiveIterator, len(a.aggregateMap))
	//idx := 0
	//for _, agg := range a.aggregateMap {
	//	its[idx] = agg.Iterator()
	//	idx++
	//}
	//return a.segmentStartTime, newFieldIterator(a.start, its)
	return
}

// GetBlock returns series block for saving loaded data
func (a *fieldAggregator) GetBlock(idx int, fn newBlockFunc) (series.Block, bool) {
	return nil, false
}

// Aggregate aggregates the field series into current aggregator
func (a *fieldAggregator) Aggregate(it series.FieldIterator) {
	//for it.HasNext() {
	//	primitiveIt := it.Next()
	//	if primitiveIt == nil {
	//		continue
	//	}
	//	primitiveFieldID := primitiveIt.FieldID()
	//	aggregator := a.getAggregator(primitiveFieldID, primitiveIt.AggType())
	//	for primitiveIt.HasNext() {
	//		timeSlot, value := primitiveIt.Next()
	//		aggregator.Aggregate(timeSlot, value)
	//	}
	//}
}

//// getAggregator returns the primitive aggregator by primitive id and function type
//func (a *fieldAggregator) getAggregator(primitiveFieldID field.PrimitiveID, aggType field.AggType) PrimitiveAggregator {
//	key := aggKey{
//		primitiveID: primitiveFieldID,
//		aggType:     aggType,
//	}
//	agg, ok := a.aggregateMap[key]
//	if ok {
//		return agg
//	}
//	agg = NewPrimitiveAggregator(primitiveFieldID, a.selector, aggType.AggFunc())
//	a.aggregateMap[key] = agg
//	return agg
//}

// reset resets the aggregate context for reusing
func (a *fieldAggregator) reset() {
	//for _, aggregator := range a.aggregateMap {
	//	aggregator.reset()
	//}
}
