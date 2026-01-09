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
	reflect "reflect"
	"unsafe"

	jsoniter "github.com/json-iterator/go"
)

var nodeTypes = make(map[string]reflect.Type)

func RegisterNodeType(node any) {
	nodeType := reflect.TypeOf(node)
	if nodeType.Kind() == reflect.Pointer {
		nodeType = nodeType.Elem()
	}
	nodeTypes[nodeType.String()] = nodeType
}

type JSONEncoder[T any] struct{}

func (encoder *JSONEncoder[T]) IsEmpty(ptr unsafe.Pointer) bool {
	return ptr == nil
}

func (encoder *JSONEncoder[T]) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	node := *(*T)(ptr)
	nodeType := reflect.TypeOf(node)
	stream.WriteObjectStart()
	stream.WriteObjectField("@type")
	stream.WriteString(nodeType.Elem().String())
	stream.WriteMore()
	stream.WriteObjectField("@data")
	stream.WriteVal(node)
	stream.WriteObjectEnd()
}

type JSONDecoder[T any] struct{}

func (decoder *JSONDecoder[T]) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	var node T
	var hasNodeType bool
	iter.ReadObjectCB(func(iter *jsoniter.Iterator, field string) bool {
		switch field {
		case "@type":
			nodeType := iter.ReadString()
			nType, ok := nodeTypes[nodeType]
			if !ok {
				panic(fmt.Sprintf("unmasharl error, unknown field type '%s'", nodeType))
			}
			node = reflect.New(nType).Interface().(T)
			hasNodeType = true
		case "@data":
			if hasNodeType {
				iter.ReadVal(node)
			} else {
				iter.Skip() // skip if node not initialed
			}
		default:
			iter.Skip() // skip unknown field
		}
		return true
	})
	*(*T)(ptr) = node
}
