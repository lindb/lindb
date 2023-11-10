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

package replica

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestFamilyChannelSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	set := newFamilyChannelSet()
	fc, ok := set.GetFamilyChannel(1)
	assert.False(t, ok)
	assert.Nil(t, fc)
	f := NewMockFamilyChannel(ctrl)
	set.InsertFamily(1, f)

	fc, ok = set.GetFamilyChannel(1)
	assert.True(t, ok)
	assert.NotNil(t, fc)

	entries := set.Entries()
	assert.Len(t, entries, 1)

	set.InsertFamily(2, f)
	set.RemoveFamilies(map[int64]struct{}{3: {}})
	entries = set.Entries()
	assert.Len(t, entries, 2)
	set.RemoveFamilies(map[int64]struct{}{})
	entries = set.Entries()
	assert.Len(t, entries, 2)

	set.RemoveFamilies(map[int64]struct{}{1: {}})
	fc, ok = set.GetFamilyChannel(1)
	assert.False(t, ok)
	assert.Nil(t, fc)
	entries = set.Entries()
	assert.Len(t, entries, 1)
}
