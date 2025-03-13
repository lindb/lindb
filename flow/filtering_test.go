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

package flow

import (
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"
)

func TestLowSeriesIDS_Find(t *testing.T) {
	ids := roaring.BitmapOf(1, 20, 400)
	seriesIDs := NewLowSeriesIDs(ids.GetContainerAtIndex(0))
	idx, ok := seriesIDs.Find(1)
	assert.Equal(t, 0, idx)
	assert.True(t, ok)
	idx, ok = seriesIDs.Find(10)
	assert.Equal(t, -1, idx)
	assert.False(t, ok)
	idx, ok = seriesIDs.Find(20)
	assert.Equal(t, 1, idx)
	assert.True(t, ok)
}
