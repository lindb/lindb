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
	"io"
	"sync"
)

var chunkReaderPool sync.Pool

// GetChunkReader picks a cached chunk-reader from the pool
func GetChunkReader(r io.Reader) *ChunkReader {
	reader := chunkReaderPool.Get()
	if reader == nil {
		return newChunkReader(r)
	}
	cr := reader.(*ChunkReader)
	cr.Reset(r)
	return cr
}

// PutChunkReader puts chunk-reader back to the pool
func PutChunkReader(cr *ChunkReader) {
	if cr == nil {
		return
	}
	chunkReaderPool.Put(cr)
}
