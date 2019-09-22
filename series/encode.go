package series

import (
	"github.com/lindb/lindb/pkg/stream"
)

func EncodeSeries(it Iterator) ([]byte, error) {
	writer := stream.NewBufferWriter(nil)
	for it.HasNext() {
		startTime, it := it.Next()
		if it == nil {
			continue
		}
		writer.PutVarint64(startTime)
		data, err := it.Bytes()
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
