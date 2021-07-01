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
	"github.com/lindb/lindb/pkg/option"
	"github.com/lindb/lindb/pkg/state"
)

func TestDatabaseService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	db := NewDatabaseService(context.TODO(), repo)

	database := models.Database{
		Name:          "test",
		Cluster:       "cluster-test",
		NumOfShard:    12,
		ReplicaFactor: 3,
		Option:        option.DatabaseOption{Interval: "10s"},
	}
	data, _ := json.Marshal(&database)

	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err := db.Save(&database)
	if err != nil {
		t.Fatal(err)
	}

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(data, nil)
	database2, _ := db.Get("test")
	assert.Equal(t, database, *database2)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
	database2, err = db.Get("test_not_exist")
	assert.Equal(t, state.ErrNotExist, err)
	assert.Nil(t, database2)
	database2, err = db.Get("")
	assert.NotNil(t, err)
	assert.Nil(t, database2)

	// json unmarshal error
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte{1, 1, 1}, nil)
	database2, err = db.Get("json_unmarshal_err")
	assert.NotNil(t, err)
	assert.Nil(t, database2)

	// test create database error
	err = db.Save(&models.Database{})
	assert.NotNil(t, err)

	err = db.Save(&models.Database{Name: "test"})
	assert.NotNil(t, err)

	err = db.Save(&models.Database{
		Name:          "test",
		Cluster:       "cluster-test",
		NumOfShard:    12,
		ReplicaFactor: 3,
	})
	assert.NotNil(t, err)

	err = db.Save(&models.Database{
		Name:          "test",
		Cluster:       "cluster-test",
		ReplicaFactor: 3,
	})
	assert.NotNil(t, err)

	err = db.Save(&models.Database{
		Name:       "test",
		Cluster:    "cluster-test",
		NumOfShard: 3,
	})
	assert.NotNil(t, err)
}

func TestDatabaseService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	db := NewDatabaseService(context.TODO(), repo)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	list, err := db.List()
	assert.NotNil(t, err)
	assert.Nil(t, list)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	list, err = db.List()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(list))
	database := models.Database{
		Name:          "test",
		Cluster:       "cluster-test",
		NumOfShard:    12,
		ReplicaFactor: 3,
	}
	database.Desc = database.String()
	data, _ := json.Marshal(&database)
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Key: "db", Value: data},
		{Key: "err", Value: []byte{1, 2, 4}},
	}, nil)
	list, err = db.List()
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 1, len(list))
	assert.Equal(t, database, *(list[0]))
}
