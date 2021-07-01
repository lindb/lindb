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

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/state"
)

func TestStorageClusterService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	repo := state.NewMockRepository(ctrl)

	cluster := config.StorageCluster{
		Name: "test1",
	}
	srv := NewStorageClusterService(context.TODO(), repo)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	err := srv.Save(&cluster)
	if err != nil {
		t.Fatal(err)
	}

	err = srv.Save(&config.StorageCluster{})
	assert.NotNil(t, err)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = srv.Save(&cluster)
	assert.NotNil(t, err)

	data, _ := json.Marshal(cluster)
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(data, nil)
	cluster2, _ := srv.Get("test1")
	assert.Equal(t, cluster, *cluster2)

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte{1, 2, 3}, nil)
	cluster2, err = srv.Get("test1")
	assert.NotNil(t, err)
	assert.Nil(t, cluster2)

	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	cluster2, err = srv.Get("test1_err")
	assert.NotNil(t, err)
	assert.Nil(t, cluster2)

	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	_ = srv.Save(&config.StorageCluster{
		Name: "test2",
	})

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	_, err1 := srv.List()
	assert.NotNil(t, err1)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Key: "data1", Value: data},
		{Key: "data2", Value: data},
		{Key: "data_err", Value: []byte{1, 2, 2}},
	}, nil)
	clusterList, _ := srv.List()
	assert.Equal(t, 2, len(clusterList))

	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
	_ = srv.Delete("test1")
}
