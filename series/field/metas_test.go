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
	"sort"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Metas(t *testing.T) {
	var metas = Metas{}
	ids := make(map[uint16]struct{})
	for i := uint16(0); i < 230; i++ {
		ids[i] = struct{}{}
	}

	for i := range ids {
		metas = metas.Insert(Meta{ID: ID(i), Type: SumField, Name: Name(strconv.Itoa(int(i)))})
	}
	sort.Sort(metas)

	// GetFromName
	m, ok := metas.GetFromName("203")
	assert.True(t, ok)
	assert.Equal(t, ID(203), m.ID)
	clone := metas.Clone()
	m, ok = clone.GetFromName("204")
	assert.True(t, ok)
	assert.Equal(t, ID(204), m.ID)
	_, ok = metas.GetFromName("250")
	assert.False(t, ok)

	// GetFromID
	m, ok = metas.GetFromID(ID(204))
	assert.True(t, ok)
	assert.Equal(t, ID(204), m.ID)
	_, ok = metas.GetFromID(ID(230))
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
	metas = metas.Insert(Meta{ID: 1, Name: "a"})
	assert.Equal(t, "a", metas.String())
	metas = metas.Insert(Meta{ID: 2, Name: "b"})
	assert.Equal(t, "a,b", metas.String())
	metas = metas.Insert(Meta{ID: 3, Name: "c,"})
	assert.Equal(t, "a,b,c,", metas.String())

}
