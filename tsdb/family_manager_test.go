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

package tsdb

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestFamilyManager_AddFamily(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	shard1 := NewMockShard(ctrl)
	shard1.EXPECT().Indicator().Return("shard1").AnyTimes()
	family1 := NewMockDataFamily(ctrl)
	family1.EXPECT().Indicator().Return("family1").AnyTimes()
	family1.EXPECT().Shard().Return(shard1)
	GetFamilyManager().AddFamily(family1)
	GetFamilyManager().AddFamily(family1)
	shard2 := NewMockShard(ctrl)
	shard2.EXPECT().Indicator().Return("shard2")
	family2 := NewMockDataFamily(ctrl)
	family2.EXPECT().Indicator().Return("family2").AnyTimes()
	family2.EXPECT().Shard().Return(shard2)
	GetFamilyManager().AddFamily(family2)
	defer GetFamilyManager().RemoveFamily(family2)
	families := GetFamilyManager().GetFamiliesByShard(shard1)
	assert.Len(t, families, 1)

	c := 0
	GetFamilyManager().WalkEntry(func(family DataFamily) {
		c++
	})
	assert.Equal(t, 2, c)
	GetFamilyManager().RemoveFamily(family1)

	c = 0
	GetFamilyManager().WalkEntry(func(family DataFamily) {
		c++
	})
	assert.Equal(t, 1, c)
}
