package series

import (
	"github.com/lindb/lindb/pkg/stream"
)

func EncodeSeries(it Iterator) ([]byte, error) {
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
		data, err := fIt.Bytes()
		if err != nil {
			return nil, err
		}
		length := len(data)
		writer.PutVarint32(int32(length))
		if length > 0 {
			writer.PutBytes(data)
		}
	}
	return writer.Bytes()
}
