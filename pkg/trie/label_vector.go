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

package trie

import (
	"bytes"
	"fmt"
	"io"
	"sort"

	"github.com/lindb/lindb/pkg/encoding"
)

const labelTerminator = 0xff

type labelVector struct {
	labels []byte
}

func labelSize(levels []*Level, endLevel int) int {
	numBytes := 0
	for l := 0; l < endLevel; l++ {
		numBytes += len(levels[l].lsLabels)
	}
	return numBytes
}

func (v *labelVector) Init(levels []*Level, endLevel uint32) {
	numBytes := labelSize(levels, int(endLevel))
	v.labels = make([]byte, numBytes)

	var pos uint32
	for l := uint32(0); l < endLevel; l++ {
		copy(v.labels[pos:], levels[l].lsLabels)
		pos += uint32(len(levels[l].lsLabels))
	}
}

func (v *labelVector) GetLabel(pos uint32) byte {
	return v.labels[pos]
}

func (v *labelVector) Search(k byte, off, size uint32) (uint32, bool) {
	start := off
	if size > 1 && v.labels[start] == labelTerminator {
		start++
		size--
	}

	end := start + size
	if end > uint32(len(v.labels)) {
		end = uint32(len(v.labels))
	}

	result := bytes.IndexByte(v.labels[start:end], k)
	if result < 0 {
		return off, false
	}
	return start + uint32(result), true
}

func (v *labelVector) SearchGreaterThan(label byte, pos, size uint32) (uint32, bool) {
	if size > 1 && v.labels[pos] == labelTerminator {
		pos++
		size--
	}

	result := sort.Search(int(size), func(i int) bool { return v.labels[pos+uint32(i)] > label })
	if uint32(result) == size {
		return pos + uint32(result) - 1, false
	}
	return pos + uint32(result), true
}

func (v *labelVector) MarshalSize(levels []*Level) int {
	size := 4
	for _, level := range levels {
		size += len(level.lsLabels)
	}
	return size
}

func (v *labelVector) Write(w io.Writer, levels []*Level) error {
	numBytes := labelSize(levels, len(levels))
	var bs [4]byte
	endian.PutUint32(bs[:], uint32(numBytes))
	if _, err := w.Write(bs[:]); err != nil {
		return err
	}
	for level := range levels {
		levelObj := levels[level]
		if _, err := w.Write(levelObj.lsLabels); err != nil {
			return err
		}
	}
	return nil
}

func (v *labelVector) Unmarshal(buf []byte) ([]byte, error) {
	if len(buf) < 4 {
		return nil, fmt.Errorf("failed unmarsh labelVector, length: %d is too short", len(buf))
	}
	size := endian.Uint32(buf)
	v.labels = buf[4 : 4+size]
	return buf[4+size:], nil
}

type valueVector struct {
	values []uint32
}

func (v *valueVector) Init(levels []*Level) {
	size := 0
	for _, level := range levels {
		size += len(level.values)
	}
	v.values = make([]uint32, size)

	pos := 0
	for level := range levels {
		values := levels[level].values
		for _, val := range values {
			v.values[pos] = val
			pos++
		}
	}
}

func (v *valueVector) Get(pos uint32) uint32 {
	return v.values[pos]
}

func (v *valueVector) Unmarshal(totalKeys int, buf []byte) ([]byte, error) {
	dataLen := totalKeys * 4
	end := dataLen
	v.values = encoding.BytesToU32Slice(buf[:end])
	return buf[end:], nil
}

type compressPathVector struct {
	hasPathVector rankVectorSparse
	offsets       []uint32
	data          []byte
}

func (pv *compressPathVector) Init(levels []*Level, bitmapType BitmapType) {
	pv.hasPathVector.Init(levels, bitmapType)
	// reset, because trie build reuse
	pv.offsets = pv.offsets[:0]
	pv.data = pv.data[:0]
	var offset uint32
	for _, l := range levels {
		var level [][]byte
		if bitmapType == HasSuffix {
			level = l.suffixes
		} else {
			level = l.prefixes
		}

		for idx := range level {
			pv.offsets = append(pv.offsets, offset)
			offset += uint32(len(level[idx]))
			pv.data = append(pv.data, level[idx]...)
		}
	}
}

func (pv *compressPathVector) MarshalSize() int {
	return pv.hasPathVector.MarshalSize() + 8 + len(pv.offsets)*4 + len(pv.data)
}

func (pv *compressPathVector) GetPath(nodeID uint32) []byte {
	if !pv.hasPathVector.IsSet(nodeID) {
		return nil
	}
	pathID := pv.hasPathVector.Rank(nodeID) - 1
	start := pv.offsets[pathID]
	end := uint32(len(pv.data))
	if int(pathID+1) < len(pv.offsets) {
		end = pv.offsets[pathID+1]
	}
	return pv.data[start:end]
}

func (pv *compressPathVector) Write(w io.Writer) error {
	if err := pv.hasPathVector.Write(w); err != nil {
		return err
	}

	var length [8]byte
	endian.PutUint32(length[:4], uint32(len(pv.offsets)*4))
	endian.PutUint32(length[4:], uint32(len(pv.data)))

	if _, err := w.Write(length[:]); err != nil {
		return err
	}
	if _, err := w.Write(encoding.U32SliceToBytes(pv.offsets)); err != nil {
		return err
	}
	if _, err := w.Write(pv.data); err != nil {
		return err
	}
	return nil
}

func (pv *compressPathVector) Unmarshal(b []byte) ([]byte, error) {
	buf1, err := pv.hasPathVector.Unmarshal(b)
	if err != nil {
		return buf1, err
	}
	if len(buf1) < 8 {
		return nil, fmt.Errorf("cannot read offsetsLen and dataLen of compressPathVector")
	}
	offsetsLen := endian.Uint32(buf1[0:4])
	dataLen := endian.Uint32(buf1[4:8])
	buf1 = buf1[8:]
	// read offsets
	if uint32(len(buf1)) < offsetsLen+dataLen {
		return nil, fmt.Errorf("offsets+data: %d is longer than remaining buf: %d of compressPathVector", offsetsLen+dataLen, len(buf1))
	}
	pv.offsets = encoding.BytesToU32Slice(buf1[:offsetsLen])
	buf1 = buf1[offsetsLen:]
	// read data
	pv.data = buf1[:dataLen]
	buf1 = buf1[dataLen:]
	return buf1, nil
}

type prefixVector struct {
	compressPathVector
}

func (v *prefixVector) CheckPrefix(key []byte, depth, nodeID uint32) (uint32, bool) {
	prefix := v.GetPrefix(nodeID)
	if len(prefix) == 0 {
		return 0, true
	}
	if int(depth)+len(prefix) > len(key) {
		return 0, false
	}
	if !bytes.Equal(key[depth:depth+uint32(len(prefix))], prefix) {
		return 0, false
	}
	return uint32(len(prefix)), true
}

func (v *prefixVector) GetPrefix(nodeID uint32) []byte {
	return v.GetPath(nodeID)
}

type suffixVector struct {
	compressPathVector
}

func (v *suffixVector) CheckSuffix(key []byte, depth, nodeID uint32) bool {
	suffix := v.GetSuffix(nodeID)
	if depth+1 >= uint32(len(key)) {
		return len(suffix) == 0
	}
	return bytes.Equal(suffix, key[depth+1:])
}

func (v *suffixVector) GetSuffix(nodeID uint32) []byte {
	return v.GetPath(nodeID)
}
