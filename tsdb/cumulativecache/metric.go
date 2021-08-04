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

package cumulativecache

import (
	"encoding/binary"

	"github.com/lindb/lindb/pkg/fasttime"
	"github.com/lindb/lindb/pkg/stream"
	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
	"github.com/lindb/lindb/series/field"
	"github.com/lindb/lindb/tsdb/memdb"
)

// fieldID+fieldValue, uint16, float64
const (
	fieldIDBytes    = 2
	fieldValueBytes = 8
)

func decodeCumulativeFieldsInto(mp *memdb.MetricPoint, data []byte) (success bool) {
	itr := &encodedDataIterator{data: data}
	var (
		fieldIDIdx = 0
	)
	ensureFieldValueCumulative := func(fieldID field.ID, value float64) (delta float64, ok bool) {
		if !itr.HasNext() {
			return 0, false
		}
		fID, fValue := itr.Next()
		if fValue > value {
			return 0, false
		}
		if fID != fieldID {
			return 0, false
		}
		return value - fValue, true
	}

	simpleFields := mp.Proto.SimpleFields
	for sfIdx := range simpleFields {
		if simpleFields[sfIdx].Type == protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM {
			delta, ok := ensureFieldValueCumulative(mp.FieldIDs[fieldIDIdx], simpleFields[sfIdx].Value)
			if !ok {
				return false
			}
			// replace field value
			simpleFields[sfIdx].Value = delta
		}
		fieldIDIdx++
	}
	// already updated
	if mp.Proto.CompoundField == nil || mp.Proto.CompoundField.Type !=
		protoMetricsV1.CompoundFieldType_CUMULATIVE_HISTOGRAM {
		return true
	}
	compoundField := mp.Proto.CompoundField
	// replace sum field
	delta, ok := ensureFieldValueCumulative(mp.FieldIDs[fieldIDIdx], compoundField.Sum)
	if !ok {
		return false
	}
	compoundField.Sum = delta
	fieldIDIdx++
	// replace count field
	delta, ok = ensureFieldValueCumulative(mp.FieldIDs[fieldIDIdx], compoundField.Count)
	if !ok {
		return false
	}
	compoundField.Count = delta
	fieldIDIdx++
	// replace values
	for cfIdx := range compoundField.Values {
		delta, ok = ensureFieldValueCumulative(mp.FieldIDs[fieldIDIdx], compoundField.Values[cfIdx])
		if !ok {
			return false
		}
		fieldIDIdx++
		compoundField.Values[cfIdx] = delta
	}
	return true
}

type encodedDataIterator struct {
	data   []byte
	offset int
	id     field.ID
	value  float64
}

func (itr *encodedDataIterator) HasNext() bool {
	if itr.offset+fieldIDBytes+fieldValueBytes > len(itr.data) {
		return false
	}
	itr.id = field.ID(stream.ReadUint16(itr.data, itr.offset))
	itr.offset += fieldIDBytes
	itr.value = float64(stream.ReadUint64(itr.data, itr.offset))
	itr.offset += fieldValueBytes
	return true
}

func (itr *encodedDataIterator) Next() (fID field.ID, v float64) {
	return itr.id, itr.value
}

func encodeCumulativeFields(mp *memdb.MetricPoint) []byte {
	var (
		fieldIDIdx       = 0
		cumulativeFields = make([]byte, timestampSizeInBytes+countCumulativeFields(mp)*(fieldIDBytes+fieldValueBytes))
		writeOffset      = 4
	)
	binary.LittleEndian.PutUint32(cumulativeFields, uint32(fasttime.UnixTimestamp()))

	putField := func(value float64) {
		stream.PutUint16(cumulativeFields, writeOffset, uint16(mp.FieldIDs[fieldIDIdx]))
		writeOffset += fieldIDBytes
		stream.PutUint64(cumulativeFields, writeOffset, uint64(value))
		writeOffset += fieldValueBytes
	}

	simpleFields := mp.Proto.SimpleFields
	for sfIdx := range simpleFields {
		if simpleFields[sfIdx].Type == protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM {
			putField(simpleFields[sfIdx].Value)
		}
		fieldIDIdx++
	}
	compoundField := mp.Proto.CompoundField
	if compoundField == nil {
		return cumulativeFields
	}
	if compoundField.Type != protoMetricsV1.CompoundFieldType_CUMULATIVE_HISTOGRAM {
		return cumulativeFields
	}
	putField(compoundField.Sum)
	fieldIDIdx++
	putField(compoundField.Count)
	fieldIDIdx++

	for idx := range compoundField.Values {
		putField(compoundField.Values[idx])
		fieldIDIdx++
	}
	return cumulativeFields
}

func countCumulativeFields(mp *memdb.MetricPoint) int {
	var fieldCount = 0
	simpleFields := mp.Proto.SimpleFields
	for sfIdx := range simpleFields {
		if simpleFields[sfIdx].Type == protoMetricsV1.SimpleFieldType_CUMULATIVE_SUM {
			fieldCount++
		}
	}
	compoundField := mp.Proto.CompoundField
	if compoundField == nil {
		return fieldCount
	}
	if compoundField.Type != protoMetricsV1.CompoundFieldType_CUMULATIVE_HISTOGRAM {
		return fieldCount
	}
	// cumulative compound field like prometheus's histogram does not contains max/min field
	return fieldCount + len(compoundField.Values) + 2
}
