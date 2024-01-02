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
	if nodeType.Kind() == reflect.Ptr {
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
