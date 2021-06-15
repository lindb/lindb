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

package series

import (
	"fmt"
	"math"

	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/stream"
	"github.com/lindb/lindb/series/field"
)

//////////////////////////////////////////////////////
// binaryGroupedIterator implements GroupedIterator
//////////////////////////////////////////////////////
type binaryGroupedIterator struct {
	tags       string
	fields     map[field.Name][]byte
	fieldNames []field.Name

	it *BinaryIterator

	idx int
}

func NewGroupedIterator(tags string, fields map[field.Name][]byte) GroupedIterator {
	it := &binaryGroupedIterator{tags: tags, fields: fields}
	for fieldName := range fields {
		it.fieldNames = append(it.fieldNames, fieldName)
	}
	return it
}

func (g *binaryGroupedIterator) Tags() string {
	return g.tags
}

func (g *binaryGroupedIterator) HasNext() bool {
	if g.idx >= len(g.fieldNames) {
		return false
	}
	g.idx++
	return true
}

func (g *binaryGroupedIterator) Next() Iterator {
	fieldName := g.fieldNames[g.idx-1]
	if g.it == nil {
		g.it = NewIterator(fieldName, g.fields[fieldName])
	} else {
		g.it.Reset(fieldName, g.fields[fieldName])
	}
	return g.it
}

//////////////////////////////////////////////////////
// BinaryIterator implements Iterator
//////////////////////////////////////////////////////
type BinaryIterator struct {
	fieldName field.Name
	fieldType field.Type
	reader    *stream.Reader
	fieldIt   *BinaryFieldIterator
	data      []byte
}

func NewIterator(fieldName field.Name, data []byte) *BinaryIterator {
	it := &BinaryIterator{fieldName: fieldName, reader: stream.NewReader(data), data: data}
	it.fieldType = field.Type(it.reader.ReadByte())
	return it
}

func (b *BinaryIterator) Reset(fieldName field.Name, data []byte) {
	b.fieldName = fieldName
	b.reader.Reset(data)
	b.fieldType = field.Type(b.reader.ReadByte())
}

func (b *BinaryIterator) FieldName() field.Name {
	return b.fieldName
}

func (b *BinaryIterator) FieldType() field.Type {
	return b.fieldType
}

func (b *BinaryIterator) HasNext() bool {
	return !b.reader.Empty()
}

func (b *BinaryIterator) Next() (startTime int64, fieldIt FieldIterator) {
	startTime = b.reader.ReadVarint64()
	length := b.reader.ReadVarint32()
	if length == 0 {
		return
	}
	data := b.reader.ReadBytes(int(length))
	if b.fieldIt == nil {
		b.fieldIt = NewFieldIterator(data)
	} else {
		b.fieldIt.reset(data)
	}
	fieldIt = b.fieldIt
	return
}

func (b *BinaryIterator) MarshalBinary() ([]byte, error) {
	return b.data, nil
}

//////////////////////////////////////////////////////
// binaryFieldIterator implements FieldIterator
//////////////////////////////////////////////////////
type BinaryFieldIterator struct {
	reader *stream.Reader
	pIt    *BinaryPrimitiveIterator
}

// NewFieldIterator create field iterator based on binary data
func NewFieldIterator(data []byte) *BinaryFieldIterator {
	it := &BinaryFieldIterator{
		reader: stream.NewReader(data),
	}
	return it
}

func (it *BinaryFieldIterator) reset(data []byte) {
	it.reader.Reset(data)
}

func (it *BinaryFieldIterator) HasNext() bool { return !it.reader.Empty() }

func (it *BinaryFieldIterator) Next() PrimitiveIterator {
	aggType := field.AggType(it.reader.ReadByte())
	length := it.reader.ReadVarint32()
	data := it.reader.ReadBytes(int(length))

	if it.pIt == nil {
		it.pIt = NewPrimitiveIterator(aggType, encoding.NewTSDDecoder(data)) //TODO get from pool?
	} else {
		it.pIt.Reset(aggType, data)
	}
	return it.pIt
}

func (it *BinaryFieldIterator) MarshalBinary() ([]byte, error) {
	return nil, fmt.Errorf("not support")
}

//////////////////////////////////////////////////////
// BinaryPrimitiveIterator implements PrimitiveIterator
//////////////////////////////////////////////////////
type BinaryPrimitiveIterator struct {
	aggType field.AggType
	tsd     *encoding.TSDDecoder
}

func NewPrimitiveIterator(aggType field.AggType, tsd *encoding.TSDDecoder) *BinaryPrimitiveIterator {
	return &BinaryPrimitiveIterator{
		aggType: aggType,
		tsd:     tsd,
	}
}

func (pi *BinaryPrimitiveIterator) Reset(aggType field.AggType, data []byte) {
	pi.aggType = aggType
	pi.tsd.Reset(data)
}

func (pi *BinaryPrimitiveIterator) AggType() field.AggType {
	return pi.aggType
}

func (pi *BinaryPrimitiveIterator) HasNext() bool {
	if pi.tsd.Error() != nil {
		return false
	}
	for pi.tsd.Next() {
		if pi.tsd.HasValue() {
			return true
		}
	}
	return false
}

func (pi *BinaryPrimitiveIterator) Next() (timeSlot int, value float64) {
	//FIXME
	timeSlot = int(pi.tsd.Slot())
	val := pi.tsd.Value()
	value = math.Float64frombits(val)
	return
}
