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
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/tsdb"
)

var testPath = "test_data"
var validOption = option.DatabaseOption{Interval: "10s"}

func TestStorageService_CreateShards(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	mockEngine := tsdb.NewMockEngine(ctrl)
	mockDatabase := tsdb.NewMockDatabase(ctrl)

	service := NewStorageService(mockEngine)
	// 2 times for double check
	err := service.CreateShards("test_db", option.DatabaseOption{})
	assert.NotNil(t, err)

	mockEngine.EXPECT().GetDatabase(gomock.Any()).Return(nil, false).MaxTimes(2)
	mockEngine.EXPECT().CreateDatabase(gomock.Any()).Return(mockDatabase, nil)
	mockDatabase.EXPECT().CreateShards(gomock.Any(), gomock.Any()).Return(nil)
	err = service.CreateShards("test_db", validOption, 1, 2, 3)
	assert.Nil(t, err)

	shard := tsdb.NewMockShard(ctrl)
	mockEngine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
	mockDatabase.EXPECT().GetShard(int32(1)).Return(shard, true)
	_, ok := service.GetShard("test_db", 1)
	assert.True(t, ok)

	mockEngine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
	mockDatabase.EXPECT().GetShard(int32(10)).Return(nil, false)
	_, ok = service.GetShard("test_db", 10)
	assert.False(t, ok)

	mockEngine.EXPECT().GetDatabase(gomock.Any()).Return(nil, false)
	_, ok = service.GetShard("not_exist_db", 10)
	assert.False(t, ok)

	// create engine error
	mockEngine.EXPECT().GetDatabase(gomock.Any()).Return(nil, false).MaxTimes(2)
	mockEngine.EXPECT().CreateDatabase("database_err").Return(nil, fmt.Errorf("err"))
	err = service.CreateShards("database_err", validOption, 1, 2, 3)
	assert.NotNil(t, err)

	// create shard error
	mockEngine.EXPECT().GetDatabase(gomock.Any()).Return(mockDatabase, true)
	mockDatabase.EXPECT().CreateShards(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = service.CreateShards("test_db", validOption, 5)
	assert.NotNil(t, err)
}

func TestStorageService_FlushDatabase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	mockEngine := tsdb.NewMockEngine(ctrl)

	service := NewStorageService(mockEngine)
	mockEngine.EXPECT().FlushDatabase(gomock.Any(), gomock.Any()).Return(true)
	ok := service.FlushDatabase(context.TODO(), "db")
	assert.True(t, ok)

}

func TestStorageService_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	mockEngine := tsdb.NewMockEngine(ctrl)
	service := NewStorageService(mockEngine)
	mockEngine.EXPECT().Close()
	service.Close()
}
