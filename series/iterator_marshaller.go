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
	"github.com/lindb/lindb/pkg/stream"
)

// MarshalIterator represents marshal series data of one field.
// format: 1byte(field type) + vint64(start time) + vint32(data length) + data
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
		length := len(data)
		writer.PutVarint32(int32(length))
		if length > 0 {
			writer.PutBytes(data)
		}
	}
	return writer.Bytes()
}
