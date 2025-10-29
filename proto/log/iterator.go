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

package log

import (
	"github.com/lindb/common/proto/gen/v1/flatLogV1"
)

type FieldIterator struct {
	l     *flatLogV1.Log
	field flatLogV1.Field

	idx int
	num int
}

func NewFieldIterator(l *flatLogV1.Log) *FieldIterator {
	return &FieldIterator{
		l:   l,
		idx: -1,
		num: l.FieldsLength(),
	}
}

func (it *FieldIterator) HasNext() bool {
	it.idx++
	if it.idx >= it.num {
		return false
	}
	return it.l.Fields(&it.field, it.idx)
}

func (it *FieldIterator) NextName() []byte { return it.field.Name() }

func (it *FieldIterator) NextValue() []byte { return it.field.Value() }

func (it *FieldIterator) Len() int { return it.num }

func (it *FieldIterator) Reset() { it.idx = -1 }
