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

package common

import (
	"bytes"
	"fmt"
	"io"

	"github.com/lindb/lindb/constants"
)

const defaultPayloadBlock = 64 * 1024

// ChunkReader reads multi lines delimited by '\n' to prevent using ioutil.ReadAll, it implements
//
// HasNext() bool (skip empty lines)
// Next() []byte
// Error() error
// Reset()
type ChunkReader struct {
	blockSize    int
	reader       io.Reader
	payloadBlock []byte
	readAt       int
	endAt        int
	section      []byte
	error        error
}

func newChunkReader(r io.Reader) *ChunkReader {
	return newChunkReaderWithSize(r, defaultPayloadBlock)
}

func newChunkReaderWithSize(r io.Reader, blockSize int) *ChunkReader {
	return &ChunkReader{
		reader:       r,
		blockSize:    blockSize,
		payloadBlock: make([]byte, blockSize),
		readAt:       0,
		endAt:        0,
	}
}

func (cr *ChunkReader) Reset(r io.Reader) {
	cr.reader = r
	cr.readAt, cr.endAt = 0, 0
	cr.section = nil
	cr.error = nil
}

func (cr *ChunkReader) moveTailToHead() {
	// copy unread data from tail to head
	unreadLen := cr.endAt - cr.readAt
	for i := 0; i < unreadLen; i++ {
		cr.payloadBlock[i] = cr.payloadBlock[i+cr.readAt]
	}
	// move cursor to the head
	cr.readAt, cr.endAt = 0, unreadLen
}

func (cr *ChunkReader) HasNext() bool {
	for {
		delimiterAt := bytes.IndexByte(cr.payloadBlock[cr.readAt:cr.endAt], '\n')
		switch delimiterAt {
		case 0: // empty line, move right
			cr.readAt++
			continue
		case -1: // do not exist
			cr.moveTailToHead()

			// re-read from reader
			n, err := cr.reader.Read(cr.payloadBlock[cr.endAt:])
			cr.error = err
			if n == 0 { // cannot read any more data
				// line is too long
				cr.section = cr.payloadBlock[cr.readAt:cr.endAt]
				if cr.endAt > len(cr.payloadBlock)-1 {
					cr.error = fmt.Errorf("%w, length %v", constants.ErrInfluxLineTooLong, cr.blockSize)
					return false
				}
				// got a line, but without delimiter
				if cr.readAt < cr.endAt {
					cr.readAt = cr.endAt
					return true
					// exhausted
				} else if cr.readAt >= cr.endAt {
					return false
				}
			} else {
				cr.endAt += n
				continue
			}
		default: // got line
			cr.section = cr.payloadBlock[cr.readAt : cr.readAt+delimiterAt]
			cr.readAt += delimiterAt + 1
			return true
		}
	}
}

func (cr *ChunkReader) Next() []byte {
	// strip prefix whitespace
	for len(cr.section) > 0 && cr.section[0] == ' ' {
		cr.section = cr.section[1:]
	}
	return cr.section
}

func (cr *ChunkReader) Error() error {
	return cr.error
}
