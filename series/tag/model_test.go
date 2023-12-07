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

import (
	"bytes"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/mock"
)

func Test_Tags(t *testing.T) {
	var tags = Tags{}
	assert.Len(t, tags.AppendHashKey(nil), 0)
	tags = append(tags, NewTag([]byte("ip"), []byte("1.1.1.1")),
		NewTag([]byte("zone"), []byte("sh")),
		NewTag([]byte("host"), []byte("test")))
	assert.Equal(t, 23, tags.Size())
	assert.False(t, tags.needsEscape())
	assert.Equal(t, ",ip=1.1.1.1,zone=sh,host=test", tags.String())
	tags = append(tags, NewTag([]byte("x x"), []byte("y,y")))
	sort.Sort(tags)
	assert.Equal(t, tags, tags.Clone())
	assert.True(t, tags.needsEscape())
	assert.Equal(t, ",host=test,ip=1.1.1.1,x\\ x=y\\,y,zone=sh", tags.String())

	tags = Tags{NewTag([]byte("x x"), []byte("y,y"))}
	m := tags.Map()
	assert.Equal(t, map[string]string{"x x": "y,y"}, m)
	assert.Len(t, TagsFromMap(m), 1)
}

func TestMeta(t *testing.T) {
	metas := Metas{{Key: "key2"}, {Key: "key1"}}
	sort.Sort(metas)
	assert.Equal(t, Metas{{Key: "key1"}, {Key: "key2"}}, metas)
	m, ok := metas.Find("key1")
	assert.True(t, ok)
	assert.Equal(t, metas[0], m)
	_, ok = metas.Find("key11")
	assert.False(t, ok)
}

func TestMeta_Marshal(t *testing.T) {
	meta := Meta{Key: "key2"}
	data, err := meta.MarshalBinary()
	assert.NoError(t, err)
	buf := bytes.NewBuffer([]byte{})
	assert.NoError(t, meta.Write(buf))
	data2 := buf.Bytes()
	assert.Equal(t, data, data2)
	ms, err := UnmarshalBinary(data)
	assert.NoError(t, err)
	assert.Equal(t, ms[0], meta)
	m2 := &Meta{}
	left := m2.Unmarshal(data)
	assert.Empty(t, left)
	assert.Equal(t, &meta, m2)
}

func TestMeta_Write_Error(t *testing.T) {
	meta := &Meta{Key: "key2"}
	assert.Error(t, meta.Write(mock.NewWriter(1)))
	assert.Error(t, meta.Write(mock.NewWriter(2)))
	assert.Error(t, meta.Write(mock.NewWriter(3)))
}
