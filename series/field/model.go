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

package field

import (
	"bytes"
	"sort"
)

type Field struct {
	Name  []byte
	Type  Type
	Value interface{}
}

// Fields implements sort.Interface
type Fields []Field

func (fs Fields) Len() int { return len(fs) }

func (fs Fields) Swap(i, j int) { fs[i], fs[j] = fs[j], fs[i] }

func (fs Fields) Less(i, j int) bool { return bytes.Compare(fs[i].Name, fs[j].Name) < 0 }

func (fs Fields) Search(name []byte) (idx int, ok bool) {
	idx = sort.Search(fs.Len(), func(i int) bool {
		return bytes.Compare(fs[i].Name, name) >= 0
	})
	if idx >= fs.Len() || !bytes.Equal(fs[idx].Name, name) {
		return -1, false
	}
	return idx, true
}

// Insert adds or replace a Field
func (fs Fields) Insert(f Field) Fields {
	if idx, ok := fs.Search(f.Name); ok {
		fs[idx] = f
		return fs
	}
	next := fs
	next = append(next, f)
	sort.Sort(next)
	return next
}
