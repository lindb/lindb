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

package grouping

type Buffer struct {
	data []uint32

	w, r int
}

func NewBuffer() *Buffer {
	return &Buffer{
		data: []uint32{},
	}
}

func (b *Buffer) Write(v uint32) {
	b.data = append(b.data, v)
	b.w++
}

func (b *Buffer) Read() uint32 {
	v := b.data[b.r]
	b.r++
	return v
}

func (b *Buffer) GetData() []uint32 {
	return b.data
}

func (b *Buffer) Reset() {
	b.data = b.data[0:0]
	b.w = 0
	b.r = 0
}

func (b *Buffer) ResetWithData(data []uint32) {
	b.data = data
	b.w = 0
	b.r = 0
}
