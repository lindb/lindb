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

package tagindex

import (
	"encoding/binary"
	"fmt"
	"sort"

	"github.com/lindb/roaring"

	"github.com/lindb/lindb/pkg/encoding"
)

const (
	indexFooterSize = 4 + // keys position
		4 + // offsets position
		4 // crc32 checksum
)

// baseReader represents the base index reader, include basic reader context
type baseReader struct {
	buf              []byte
	offsets          *encoding.FixedOffsetDecoder
	keys             *roaring.Bitmap
	tagValueBitmapAt int
	offsetsAt        int
	crc32CheckSum    uint32
}

// initReader initializes the basic index reader context
func (r *baseReader) initReader() error {
	if len(r.buf) <= indexFooterSize {
		return fmt.Errorf("block length short:%d than footer size: %d", len(r.buf), indexFooterSize)
	}
	// read footer(4+4+4)
	footerPos := len(r.buf) - indexFooterSize
	r.tagValueBitmapAt = int(binary.LittleEndian.Uint32(r.buf[footerPos : footerPos+4]))
	r.offsetsAt = int(binary.LittleEndian.Uint32(r.buf[footerPos+4 : footerPos+8]))
	r.crc32CheckSum = binary.LittleEndian.Uint32(r.buf[footerPos+8 : footerPos+12])
	// validate offsets
	if !sort.IntsAreSorted([]int{
		0, r.tagValueBitmapAt, r.offsetsAt, footerPos}) {
		return fmt.Errorf("invalid footer format")
	}
	// read keys
	keys := roaring.New()
	if err := encoding.BitmapUnmarshal(keys, r.buf[r.tagValueBitmapAt:]); err != nil {
		return err
	}
	r.keys = keys
	// read high keys offsets
	r.offsets = encoding.NewFixedOffsetDecoder()
	_, err := r.offsets.Unmarshal(r.buf[r.offsetsAt:])
	return err
}
