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

package tag

import "bytes"

// Meta holds the relation of tagKey and its ID
type Meta struct {
	Key string
	ID  uint32
}

// Tag represents a kv tag pair.
type Tag struct {
	Key   []byte
	Value []byte
}

// Size returns the slice's size of the key and value.
func (t Tag) Size() int { return len(t.Key) + len(t.Value) }

// NewTag returns a new Tag
func NewTag(key, value []byte) Tag {
	return Tag{Key: key, Value: value}
}

// Tags implements sort.Interface
type Tags []Tag

func (tags Tags) Len() int           { return len(tags) }
func (tags Tags) Swap(i, j int)      { tags[i], tags[j] = tags[j], tags[i] }
func (tags Tags) Less(i, j int) bool { return bytes.Compare(tags[i].Key, tags[j].Key) < 0 }
func (tags Tags) Size() int {
	var total int
	for i := range tags {
		total += tags[i].Size()
	}
	return total
}

func (tags Tags) String() string {
	return string(tags.AppendHashKey(nil))
}

func (tags Tags) needsEscape() bool {
	for i := range tags {
		t := &tags[i]
		for j := range tagEscapeCodes {
			c := &tagEscapeCodes[j]
			if bytes.IndexByte(t.Key, c.k[0]) != -1 || bytes.IndexByte(t.Value, c.k[0]) != -1 {
				return true
			}
		}
	}
	return false
}

// AppendHashKey appends the result of hashing all of a tag's keys and values to dst and returns the extended buffer.
func (tags Tags) AppendHashKey(dst []byte) []byte {
	// Empty maps marshal to empty bytes.
	if len(tags) == 0 {
		return dst
	}

	// Type invariant: Tags are sorted
	sz := 0
	var escaped Tags
	if tags.needsEscape() {
		var tmp [20]Tag
		if len(tags) < len(tmp) {
			escaped = tmp[:len(tags)]
		} else {
			escaped = make(Tags, len(tags))
		}

		for i := range tags {
			t := &tags[i]
			nt := &escaped[i]
			nt.Key = EscapeTag(t.Key)
			nt.Value = EscapeTag(t.Value)
			sz += len(nt.Key) + len(nt.Value)
		}
	} else {
		sz = tags.Size()
		escaped = tags
	}

	sz += len(escaped) + (len(escaped) * 2) // separators

	// Generate marshaled bytes.
	if cap(dst)-len(dst) < sz {
		nd := make([]byte, len(dst), len(dst)+sz)
		copy(nd, dst)
		dst = nd
	}
	buf := dst[len(dst) : len(dst)+sz]
	idx := 0
	for i := range escaped {
		k := &escaped[i]
		if len(k.Value) == 0 {
			continue
		}
		buf[idx] = ','
		idx++
		copy(buf[idx:], k.Key)
		idx += len(k.Key)
		buf[idx] = '='
		idx++
		copy(buf[idx:], k.Value)
		idx += len(k.Value)
	}
	return dst[:len(dst)+idx]
}
