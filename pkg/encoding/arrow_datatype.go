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
	"strconv"
	"strings"
	"unsafe"

	"github.com/apache/arrow-go/v18/arrow"
	jsoniter "github.com/json-iterator/go"
)

var arrowDataTypes = map[string]arrow.DataType{
	arrow.BinaryTypes.String.Name():           arrow.BinaryTypes.String,
	arrow.BinaryTypes.Binary.Name():           arrow.BinaryTypes.Binary,
	arrow.PrimitiveTypes.Int8.Name():          arrow.PrimitiveTypes.Int8,
	arrow.PrimitiveTypes.Int16.Name():         arrow.PrimitiveTypes.Int16,
	arrow.PrimitiveTypes.Int32.Name():         arrow.PrimitiveTypes.Int32,
	arrow.PrimitiveTypes.Int64.Name():         arrow.PrimitiveTypes.Int64,
	arrow.PrimitiveTypes.Uint8.Name():         arrow.PrimitiveTypes.Uint8,
	arrow.PrimitiveTypes.Uint16.Name():        arrow.PrimitiveTypes.Uint16,
	arrow.PrimitiveTypes.Uint32.Name():        arrow.PrimitiveTypes.Uint32,
	arrow.PrimitiveTypes.Uint64.Name():        arrow.PrimitiveTypes.Uint64,
	arrow.PrimitiveTypes.Float32.Name():       arrow.PrimitiveTypes.Float32,
	arrow.PrimitiveTypes.Float64.Name():       arrow.PrimitiveTypes.Float64,
	arrow.FixedWidthTypes.Timestamp_ns.Name(): arrow.FixedWidthTypes.Timestamp_ns,
	// FIXME: same key arrow.FixedWidthTypes.Timestamp_ms.Name():  arrow.FixedWidthTypes.Timestamp_ms,
	arrow.FixedWidthTypes.Boolean.Name():     arrow.FixedWidthTypes.Boolean,
	arrow.FixedWidthTypes.Duration_ns.Name(): arrow.FixedWidthTypes.Duration_ns,
	// FIXME: same key arrow.FixedWidthTypes.Duration_ms.Name():   arrow.FixedWidthTypes.Duration_ms,
	arrow.Null.Name(): arrow.Null,
}

// RegisterArrowDataType registers an arrow.DataType so it can be looked up
// by name during JSON deserialization. This must be called for any extension
// types (e.g. TimeSeries, Exemplar, Aggregation) before they are decoded.
func RegisterArrowDataType(dt arrow.DataType) {
	arrowDataTypes[dt.Name()] = dt
}

// encodeArrowType encodes an Arrow DataType to its canonical string representation.
// Parametric and nested types are encoded recursively so the full structure is preserved.
//
// Encoding rules:
//
//	FixedSizeBinary(N)       → "fixed_size_binary[N]"
//	Map<K,V>                 → "map<{encode(K)},{encode(V)}>"
//	List<T>                  → "list<{encode(T)}>"
//	Struct<f1:T1,f2:T2,...>  → "struct<f1={encode(T1)},f2={encode(T2)},...>"
//	anything else            → dt.Name()
func encodeArrowType(dt arrow.DataType) string {
	if fsb, ok := dt.(*arrow.FixedSizeBinaryType); ok {
		return fmt.Sprintf("fixed_size_binary[%d]", fsb.ByteWidth)
	}
	if m, ok := dt.(*arrow.MapType); ok {
		return "map<" + encodeArrowType(m.KeyType()) + "," + encodeArrowType(m.ItemType()) + ">"
	}
	if lt, ok := dt.(*arrow.ListType); ok {
		return "list<" + encodeArrowType(lt.Elem()) + ">"
	}
	if st, ok := dt.(*arrow.StructType); ok {
		var sb strings.Builder
		sb.WriteString("struct<")
		for i := 0; i < st.NumFields(); i++ {
			if i > 0 {
				sb.WriteByte(',')
			}
			f := st.Field(i)
			sb.WriteString(f.Name)
			sb.WriteByte('=')
			sb.WriteString(encodeArrowType(f.Type))
		}
		sb.WriteByte('>')
		return sb.String()
	}
	return dt.Name()
}

// decodeArrowType decodes a type string produced by encodeArrowType back into an Arrow DataType.
// Returns (nil, false) when the string is not recognised.
func decodeArrowType(name string) (arrow.DataType, bool) {
	// FixedSizeBinary[N]
	if strings.HasPrefix(name, "fixed_size_binary[") && strings.HasSuffix(name, "]") {
		inner := name[len("fixed_size_binary[") : len(name)-1]
		byteWidth, err := strconv.Atoi(inner)
		if err != nil {
			return nil, false
		}
		return &arrow.FixedSizeBinaryType{ByteWidth: byteWidth}, true
	}

	// Map<K,V>
	if strings.HasPrefix(name, "map<") && strings.HasSuffix(name, ">") {
		inner := name[4 : len(name)-1]
		parts := splitTopLevel(inner)
		if len(parts) == 2 {
			kt, ok1 := decodeArrowType(parts[0])
			vt, ok2 := decodeArrowType(parts[1])
			if ok1 && ok2 {
				return arrow.MapOf(kt, vt), true
			}
		}
		return nil, false
	}

	// List<T>
	if strings.HasPrefix(name, "list<") && strings.HasSuffix(name, ">") {
		inner := name[5 : len(name)-1]
		itemType, ok := decodeArrowType(inner)
		if ok {
			return arrow.ListOf(itemType), true
		}
		return nil, false
	}

	// Struct<f1=T1,f2=T2,...>
	if strings.HasPrefix(name, "struct<") && strings.HasSuffix(name, ">") {
		inner := name[7 : len(name)-1]
		if inner == "" {
			return arrow.StructOf(), true
		}
		fieldDefs := splitTopLevel(inner)
		fields := make([]arrow.Field, 0, len(fieldDefs))
		for _, fd := range fieldDefs {
			// Field format: "name=type"
			eqIdx := strings.IndexByte(fd, '=')
			if eqIdx < 0 {
				return nil, false
			}
			fName := fd[:eqIdx]
			fType, ok := decodeArrowType(fd[eqIdx+1:])
			if !ok {
				return nil, false
			}
			fields = append(fields, arrow.Field{Name: fName, Type: fType})
		}
		return arrow.StructOf(fields...), true
	}

	// Known simple types
	if dt, ok := arrowDataTypes[name]; ok {
		return dt, true
	}
	return nil, false
}

// splitTopLevel splits s on commas that are NOT nested inside angle/square brackets.
// This is necessary for parsing nested type strings like "struct<f1=list<utf8>,f2=utf8>".
func splitTopLevel(s string) []string {
	var parts []string
	depth := 0
	start := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '<', '[':
			depth++
		case '>', ']':
			depth--
		case ',':
			if depth == 0 {
				parts = append(parts, s[start:i])
				start = i + 1
			}
		}
	}
	parts = append(parts, s[start:])
	return parts
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
	stream.WriteString(encodeArrowType(dt))
}

type arrowDataTypeDecoder struct{}

func (arrowDataTypeDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	name := iter.ReadString()

	dt, ok := decodeArrowType(name)
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
