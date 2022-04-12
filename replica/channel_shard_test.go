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
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func TestShardChannel_SyncShardState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	ch := newShardChannel(context.TODO(), "database", 1, nil)

	familyCh := NewMockFamilyChannel(ctrl)
	ch1 := ch.(*shardChannel)
	ch1.mutex.Lock()
	ch1.families.InsertFamily(1, familyCh)
	ch1.shardState = models.ShardState{
		Leader: 1,
	}
	ch1.mutex.Unlock()

	// leader no change
	ch.SyncShardState(models.ShardState{
		Leader: 1,
	}, nil)

	// leader change
	familyCh.EXPECT().leaderChanged(gomock.Any(), gomock.Any())
	ch.SyncShardState(models.ShardState{
		Leader: 2,
	}, nil)

	ch1.mutex.Lock()
	assert.Equal(t, models.ShardState{
		Leader: 2,
	}, ch1.shardState)
	ch1.mutex.Unlock()
}

func TestShardChannel_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
	}()
	ch := newShardChannel(context.TODO(), "database", 1, nil)

	familyCh := NewMockFamilyChannel(ctrl)
	ch1 := ch.(*shardChannel)
	ch1.mutex.Lock()
	ch1.families.InsertFamily(1, familyCh)
	ch1.shardState = models.ShardState{
		Leader: 1,
	}
	ch1.mutex.Unlock()

	familyCh.EXPECT().FamilyTime().Return(int64(1)).AnyTimes()
	familyCh.EXPECT().isExpire(gomock.Any(), gomock.Any()).Return(true)
	familyCh.EXPECT().Stop(gomock.Any())
	ch.garbageCollect(1, 1)

	// no family need stop
	ch.Stop()

	// add new family for stop
	ch1.mutex.Lock()
	ch1.families.InsertFamily(1, familyCh)
	ch1.shardState = models.ShardState{
		Leader: 1,
	}
	ch1.mutex.Unlock()
	familyCh.EXPECT().Stop(gomock.Any())
	ch.Stop()
}

func TestShardChannel_GetOrCreateFamilyChannel(t *testing.T) {
	defer func() {
		getFamilyFn = getFamily
	}()
	ch := newShardChannel(context.TODO(), "database", 1, nil)
	f1 := ch.GetOrCreateFamilyChannel(1)
	assert.NotNil(t, f1)
	f2 := ch.GetOrCreateFamilyChannel(1)
	assert.Equal(t, f1, f2)

	// test double check
	getFamilyFn = func(families *familyChannelSet, familyTime int64) (FamilyChannel, bool) {
		return f1, true
	}
	f3 := ch.GetOrCreateFamilyChannel(3)
	assert.Equal(t, f1, f3)
}
