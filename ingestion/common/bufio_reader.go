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
	"bufio"
	"io"
	"sync"
)

var bufioReaderPool sync.Pool

// NewBufioReader returns a buf io reader with 64KB cache
func NewBufioReader(r io.Reader) (*bufio.Reader, func(bufioReader *bufio.Reader)) {
	putbackFunc := func(bufioReader *bufio.Reader) {
		bufioReader.Reset(nil)
		bufioReaderPool.Put(bufioReader)
	}
	var reader *bufio.Reader
	item := bufioReaderPool.Get()
	if item != nil {
		reader = item.(*bufio.Reader)
	} else {
		reader = bufio.NewReaderSize(r, 64*1024)
	}
	reader.Reset(r)
	return reader, putbackFunc
}
