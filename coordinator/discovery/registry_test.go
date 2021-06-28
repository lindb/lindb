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

package discovery

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

var testRegistryPath = "/test/registry"

func TestRegistry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	registry1 := NewRegistry(repo, testRegistryPath, 100)

	closedCh := make(chan state.Closed)

	node := models.Node{IP: "127.0.0.1", Port: 2080, HTTPPort: 9002}
	gomock.InOrder(
		repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("err")),
		repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(closedCh, nil),
	)
	err := registry1.Register(node)
	assert.NoError(t, err)
	time.Sleep(600 * time.Millisecond)

	// maybe retry do heartbeat after close chan
	repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	close(closedCh)
	time.Sleep(600 * time.Millisecond)

	nodePath := fmt.Sprintf("%s/data/%s", testRegistryPath, node.Indicator())
	repo.EXPECT().Delete(gomock.Any(), nodePath).Return(nil)
	err = registry1.Deregister(node)
	assert.Nil(t, err)

	err = registry1.Close()
	assert.NoError(t, err)

	registry1 = NewRegistry(repo, testRegistryPath, 100)
	err = registry1.Close()
	assert.NoError(t, err)

	r := registry1.(*registry)
	r.register("/data/pant", node)

	registry1 = NewRegistry(repo, testRegistryPath, 100)
	r = registry1.(*registry)

	// cancel ctx in timer
	time.AfterFunc(100*time.Millisecond, func() {
		r.cancel()
	})
	r.register("/data/pant", node)
}

func TestRegistry_GenerateNodeID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	registry1 := NewRegistry(repo, testRegistryPath, 100)
	// case 1: get err
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	_, err := registry1.GenerateNodeID(models.Node{})
	assert.Error(t, err)

	// case 2: get data ok
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return([]byte(strconv.FormatInt(10, 10)), nil)
	id, err := registry1.GenerateNodeID(models.Node{})
	assert.NoError(t, err)
	assert.Equal(t, models.NodeID(10), id)

	// case 3: init seq err
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
	repo.EXPECT().NextSequence(gomock.Any(), gomock.Any()).Return(int64(0), fmt.Errorf("err"))
	_, err = registry1.GenerateNodeID(models.Node{})
	assert.Error(t, err)
	// case 4: store seq err
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
	repo.EXPECT().NextSequence(gomock.Any(), gomock.Any()).Return(int64(10), nil)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	_, err = registry1.GenerateNodeID(models.Node{})
	assert.Error(t, err)
	// case 5: init seq
	repo.EXPECT().Get(gomock.Any(), gomock.Any()).Return(nil, state.ErrNotExist)
	repo.EXPECT().NextSequence(gomock.Any(), gomock.Any()).Return(int64(100), nil)
	repo.EXPECT().Put(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	id, err = registry1.GenerateNodeID(models.Node{})
	assert.NoError(t, err)
	assert.Equal(t, models.NodeID(100), id)
}
