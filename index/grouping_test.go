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

package index

import (
	"testing"

	"github.com/lindb/roaring"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/imap"
)

func TestGroupingScanner(t *testing.T) {
	forward := imap.NewIntMap[uint32]()
	forward.Put(1, 1)
	forward.Put(2, 2)

	scanner := &memGroupingScanner{
		forward:  forward,
		withLock: func() (release func()) { return func() {} },
	}

	assert.Equal(t, roaring.BitmapOf(1, 2).ToArray(), scanner.GetSeriesIDs().ToArray())
	lowSeriesIDs, tagValueIDs := scanner.GetSeriesAndTagValue(0)
	assert.Equal(t, []uint16{1, 2}, lowSeriesIDs.ToArray())
	assert.Equal(t, []uint32{1, 2}, tagValueIDs)
	lowSeriesIDs, tagValueIDs = scanner.GetSeriesAndTagValue(1)
	assert.Nil(t, lowSeriesIDs)
	assert.Nil(t, tagValueIDs)
}
