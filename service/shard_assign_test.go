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

package service

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestShardAssignService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	srv := NewShardAssignService(context.TODO(), repo)

	shardAssign1 := models.NewShardAssignment("test")
	shardAssign1.AddReplica(1, 1)
	shardAssign1.AddReplica(1, 2)
	shardAssign1.AddReplica(1, 3)
	shardAssign1.AddReplica(2, 2)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	_ = srv.Save("db1", shardAssign1)

	shardAssign2 := models.NewShardAssignment("test")
	shardAssign2.AddReplica(1, 1)
	shardAssign2.AddReplica(2, 2)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	_ = srv.Save("db2", shardAssign2)

	data1, _ := json.Marshal(shardAssign1)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(data1, nil)
	shardAssign11, _ := srv.Get("db1")
	assert.Equal(t, *shardAssign1, *shardAssign11)

	data2, _ := json.Marshal(shardAssign2)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(data2, nil)
	shardAssign22, _ := srv.Get("db2")
	assert.Equal(t, *shardAssign2, *shardAssign22)

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
	_, err := srv.Get("not_exist")
	assert.Equal(t, state.ErrNotExist, err)

	// unmarshal error
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte{1, 3, 34}, nil)
	_, err = srv.Get("not_exist")
	assert.NotNil(t, err)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Key: "data2", Value: data2},
		{Key: "err", Value: []byte{1, 1, 1}},
	}, nil)
	list, _ := srv.List()
	assert.Equal(t, 1, len(list))
	assert.Equal(t, *shardAssign2, *(list[0]))

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	list, err = srv.List()
	assert.Nil(t, list)
	assert.NotNil(t, err)
}
