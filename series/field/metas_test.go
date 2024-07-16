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
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/internal/mock"
)

func Test_Metas(t *testing.T) {
	metas := Metas{}
	ids := make(map[uint16]struct{})
	for i := uint16(0); i < 230; i++ {
		ids[i] = struct{}{}
	}

	for i := range ids {
		metas = append(metas, Meta{ID: ID(i), Type: SumField, Name: Name(strconv.Itoa(int(i)))})
	}
	sort.Sort(metas)

	// GetFromName
	m, ok := metas.GetFromName("203")
	assert.True(t, ok)
	assert.Equal(t, ID(203), m.ID)
	clone := metas.Clone()
	m, ok = clone.GetFromName("204")
	assert.True(t, ok)
	idx, ok := clone.FindIndexByName("1")
	assert.True(t, ok)
	assert.Equal(t, 1, idx)
	idx, ok = clone.FindIndexByName("1000")
	assert.False(t, ok)
	assert.Equal(t, -1, idx)

	assert.Equal(t, ID(204), m.ID)
	_, ok = metas.GetFromName("250")
	assert.False(t, ok)

	// GetFromID
	m, ok = metas.GetFromID(ID(204))
	assert.True(t, ok)
	assert.Equal(t, ID(204), m.ID)
	m, ok = metas.Find("204")
	assert.True(t, ok)
	assert.Equal(t, ID(204), m.ID)
	_, ok = metas.GetFromID(ID(230))
	assert.False(t, ok)
	_, ok = metas.Find("230")
	assert.False(t, ok)

	// Intersects
	ml, ok := metas.Intersects(Metas{{ID: 1}, {ID: 203}, {ID: 250}})
	assert.False(t, ok)
	assert.Len(t, ml, 2)
	ml, ok = metas.Intersects(Metas{{ID: 1}, {ID: 203}, {ID: 204}})
	assert.True(t, ok)
	assert.Len(t, ml, 3)

	metas = Metas{}
	assert.Equal(t, "", metas.String())
	metas = append(metas, Meta{ID: 1, Name: "a"})
	assert.Equal(t, "a", metas.String())
	metas = append(metas, Meta{ID: 2, Name: "b"})
	assert.Equal(t, "a,b", metas.String())
	metas = append(metas, Meta{ID: 3, Name: "c,"})
	assert.Equal(t, "a,b,c,", metas.String())
}

func TestMeta_Marshal(t *testing.T) {
	m := &Meta{
		ID:   1,
		Name: "f",
		Type: SumField,
	}
	data, err := m.MarshalBinary()
	assert.NoError(t, err)
	buf := bytes.NewBuffer([]byte{})
	assert.NoError(t, m.Write(buf))
	data2 := buf.Bytes()
	assert.Equal(t, data, data2)

	m2 := &Meta{}
	left := m2.Unmarshal(data)
	assert.Empty(t, left)
	assert.Equal(t, m2, m)
	m3, id, err := UnmarshalBinary(data)
	assert.Equal(t, *m, m3[0])
	assert.Equal(t, ID(1), id)
	assert.NoError(t, err)
}

func TestMeta_Write_Error(t *testing.T) {
	m := &Meta{}
	assert.Error(t, m.Write(mock.NewWriter(1)))
	assert.Error(t, m.Write(mock.NewWriter(2)))
	assert.Error(t, m.Write(mock.NewWriter(3)))
}

func TestMetas_Find(t *testing.T) {
	fields := Metas{
		{Name: "HistogramSum"},
		{Name: "HistogramMin"},
		{Name: "HistogramMax"},
		{Name: "HistogramCount"},
	}
	sort.Sort(fields)
	fmt.Println(fields)
	f, ok := fields.Find("HistogramMax")
	assert.True(t, ok)
	assert.Equal(t, Name("HistogramMax"), f.Name)
}
