package series

import (
	"github.com/lindb/lindb/pkg/stream"
)

func MarshalIterator(it Iterator) ([]byte, error) {
	if it == nil {
		return nil, nil
	}
	writer := stream.NewBufferWriter(nil)
	writer.PutByte(byte(it.FieldType()))
	for it.HasNext() {
		startTime, fIt := it.Next()
		if fIt == nil {
			continue
		}
		writer.PutVarint64(startTime)
		data, err := fIt.MarshalBinary()
		if err != nil {
			return nil, err
		}
		writer.PutVarint32(int32(len(data)))
		if len(data) > 0 {
			writer.PutBytes(data)
		}
	}
	return writer.Bytes()
}
