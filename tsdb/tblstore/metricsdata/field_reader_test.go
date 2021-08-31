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
	"encoding/binary"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/kv"
	"github.com/lindb/lindb/pkg/timeutil"
	"github.com/lindb/lindb/series/field"
)

func TestField_read(t *testing.T) {
	block := mockMetricMergeBlock([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner, _ := newDataScanner(r)
	seriesEntry := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), seriesEntry, timeutil.SlotRange{Start: 5, End: 5})
	sr := fReader.SlotRange()
	assert.Equal(t, uint16(5), sr.Start)
	assert.Equal(t, uint16(5), sr.End)
	// case 1: field 1 not exist
	data := fReader.GetFieldData(1)
	assert.Nil(t, data)
	// case 2: field 2 exist
	data = fReader.GetFieldData(2)
	assert.True(t, len(data) > 0)
	// case 3: field 10 exist
	data = fReader.GetFieldData(10)
	assert.True(t, len(data) > 0)
	// case 4: field 20 not exist
	data = fReader.GetFieldData(20)
	assert.Nil(t, data)
	// case 5: complete cannot get field
	fReader.Close()
	data = fReader.GetFieldData(10)
	assert.Nil(t, data)
	// case 6: no fields
	fReader = newFieldReader(scanner.fieldIndexes(), []byte{0, 0, 0}, timeutil.SlotRange{Start: 5, End: 5})
	data = fReader.GetFieldData(10)
	assert.Nil(t, data)
}

func TestFieldReader_close(t *testing.T) {
	block := mockMetricMergeBlock([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner, _ := newDataScanner(r)
	seriesEntry := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), seriesEntry, timeutil.SlotRange{Start: 5, End: 5})
	fReader.Close()
	data := fReader.GetFieldData(2)
	assert.Nil(t, data)
}

func TestFieldReader_reset(t *testing.T) {
	block := mockMetricMergeBlock([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner, _ := newDataScanner(r)
	seriesEntry := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), seriesEntry, timeutil.SlotRange{Start: 5, End: 5})
	sr := fReader.SlotRange()
	assert.Equal(t, uint16(5), sr.Start)
	assert.Equal(t, uint16(5), sr.End)
	data := fReader.GetFieldData(2)
	assert.True(t, len(data) > 0)
	data = fReader.GetFieldData(10)
	assert.True(t, len(data) > 0)

	// mock diff field
	nopKVFlusher := kv.NewNopFlusher()
	flusher, _ := NewFlusher(nopKVFlusher)
	flusher.PrepareMetric(10, field.Metas{
		{ID: 10, Type: field.MinField},
	})
	_ = flusher.FlushField([]byte{1, 2, 3})
	_ = flusher.FlushSeries(10)
	_ = flusher.CommitMetric(sr)

	// reset value
	fReader.Reset(seriesEntry, timeutil.SlotRange{Start: 15, End: 15})
	sr = fReader.SlotRange()
	assert.Equal(t, uint16(15), sr.Start)
	assert.Equal(t, uint16(15), sr.End)
	data = fReader.GetFieldData(10)
	assert.True(t, len(data) > 0)
}

func TestFieldReader_Reset_error(t *testing.T) {
	block := mockMetricMergeBlock([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner, _ := newDataScanner(r)
	seriesEntry := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), seriesEntry, timeutil.SlotRange{Start: 5, End: 5})
	fReader.Reset(nil, timeutil.SlotRange{Start: 5, End: 5})
	assert.True(t, fReader.(*fieldReader).completed)
	// max uint64
	var buf [binary.MaxVarintLen64]byte
	binary.PutUvarint(buf[:], math.MaxUint64)
	reverseBuf(buf[:])
	fReader = newFieldReader(scanner.fieldIndexes(), seriesEntry, timeutil.SlotRange{Start: 5, End: 5})
	fReader.Reset(buf[:], timeutil.SlotRange{Start: 5, End: 5})
	assert.True(t, fReader.(*fieldReader).completed)
	// bad variant
	var buf2 = []byte{
		1, 1,
		0x80, 0x80, 0x80, 0x80, 0x80,
		0x80, 0x80, 0x80, 0x80, 0x80,
	}
	reverseBuf(buf2)
	fReader = newFieldReader(scanner.fieldIndexes(), seriesEntry, timeutil.SlotRange{Start: 5, End: 5})
	fReader.Reset(buf2, timeutil.SlotRange{Start: 5, End: 5})
	assert.True(t, fReader.(*fieldReader).completed)
	// empty buf
	var buf3 = []byte{
		1, 1, 1, 1, 1,
		1, 1, 1, 1, 1,
	}
	fReader = newFieldReader(scanner.fieldIndexes(), seriesEntry, timeutil.SlotRange{Start: 5, End: 5})
	fReader.Reset(buf3, timeutil.SlotRange{Start: 5, End: 5})
	assert.True(t, fReader.(*fieldReader).completed)
}

func reverseBuf(data []byte) {
	for i := 0; i < len(data); i++ {
		data[i], data[len(data)-i-1] = data[len(data)-i-1], data[i]
	}
}

func TestFieldReader_read_one_field(t *testing.T) {
	block := mockMetricMergeBlockOneField([]uint32{1}, 5, 5)
	r, err := NewReader("1.sst", block)
	assert.NoError(t, err)
	assert.NotNil(t, r)
	scanner, _ := newDataScanner(r)
	seriesEntry := scanner.scan(0, 1)
	fReader := newFieldReader(scanner.fieldIndexes(), seriesEntry, timeutil.SlotRange{Start: 5, End: 5})
	sr := fReader.SlotRange()
	assert.Equal(t, uint16(5), sr.Start)
	assert.Equal(t, uint16(5), sr.End)
	// case 1: field 1 not exist
	data := fReader.GetFieldData(1)
	assert.Nil(t, data)
	// case 2: field 2 exist
	data = fReader.GetFieldData(2)
	assert.True(t, len(data) > 0)
	// case 3: close cannot metricReader data
	fReader.Close()
	data = fReader.GetFieldData(2)
	assert.Nil(t, data)
}
