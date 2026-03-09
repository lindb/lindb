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

package encoding

import (
	"fmt"
	"unsafe"

	"github.com/apache/arrow-go/v18/arrow"
	jsoniter "github.com/json-iterator/go"
)

var arrowDataTypes = map[string]arrow.DataType{
	arrow.BinaryTypes.String.Name():            arrow.BinaryTypes.String,
	arrow.BinaryTypes.Binary.Name():            arrow.BinaryTypes.Binary,
	arrow.PrimitiveTypes.Int8.Name():           arrow.PrimitiveTypes.Int8,
	arrow.PrimitiveTypes.Int16.Name():          arrow.PrimitiveTypes.Int16,
	arrow.PrimitiveTypes.Int32.Name():          arrow.PrimitiveTypes.Int32,
	arrow.PrimitiveTypes.Int64.Name():          arrow.PrimitiveTypes.Int64,
	arrow.PrimitiveTypes.Uint8.Name():          arrow.PrimitiveTypes.Uint8,
	arrow.PrimitiveTypes.Uint16.Name():         arrow.PrimitiveTypes.Uint16,
	arrow.PrimitiveTypes.Uint32.Name():         arrow.PrimitiveTypes.Uint32,
	arrow.PrimitiveTypes.Uint64.Name():         arrow.PrimitiveTypes.Uint64,
	arrow.PrimitiveTypes.Float32.Name():        arrow.PrimitiveTypes.Float32,
	arrow.PrimitiveTypes.Float64.Name():        arrow.PrimitiveTypes.Float64,
	arrow.FixedWidthTypes.Timestamp_ns.Name():  arrow.FixedWidthTypes.Timestamp_ns,
	arrow.FixedWidthTypes.Timestamp_ms.Name():  arrow.FixedWidthTypes.Timestamp_ms,
	arrow.FixedWidthTypes.Boolean.Name():       arrow.FixedWidthTypes.Boolean,
	arrow.FixedWidthTypes.Duration_ns.Name():   arrow.FixedWidthTypes.Duration_ns,
	arrow.FixedWidthTypes.Duration_ms.Name():   arrow.FixedWidthTypes.Duration_ms,
	arrow.Null.Name():                          arrow.Null,
}

// RegisterArrowDataType registers an arrow.DataType so it can be looked up
// by name during JSON deserialization. This must be called for any extension
// types (e.g. TimeSeries, Exemplar, Aggregation) before they are decoded.
func RegisterArrowDataType(dt arrow.DataType) {
	arrowDataTypes[dt.Name()] = dt
}

type arrowDataTypeEncoder struct{}

func (arrowDataTypeEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	return *(*arrow.DataType)(ptr) == nil
}

func (arrowDataTypeEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	dt := *(*arrow.DataType)(ptr)
	if dt == nil {
		stream.WriteNil()
		return
	}
	stream.WriteString(dt.Name())
}

type arrowDataTypeDecoder struct{}

func (arrowDataTypeDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	name := iter.ReadString()
	dt, ok := arrowDataTypes[name]
	if !ok {
		iter.ReportError("arrowDataTypeDecoder", fmt.Sprintf("unknown arrow.DataType name: %q", name))
		return
	}
	*(*arrow.DataType)(ptr) = dt
}

func init() {
	jsoniter.RegisterTypeEncoder("arrow.DataType", arrowDataTypeEncoder{})
	jsoniter.RegisterTypeDecoder("arrow.DataType", arrowDataTypeDecoder{})
}
