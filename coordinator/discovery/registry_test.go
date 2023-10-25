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
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestRegistry(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	registry1 := NewRegistry(repo, constants.LiveNodesPath, 100)

	closedCh := make(chan state.Closed)

	node := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 2080}
	gomock.InOrder(
		repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, fmt.Errorf("err")),
		repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(closedCh, nil),
	)
	err := registry1.Register(&node)
	assert.NoError(t, err)
	time.Sleep(600 * time.Millisecond)
	assert.True(t, registry1.IsSuccess())

	// maybe retry do heartbeat after close chan
	repo.EXPECT().Heartbeat(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	close(closedCh)
	time.Sleep(600 * time.Millisecond)

	nodePath := fmt.Sprintf("%s/%s", constants.LiveNodesPath, node.Indicator())
	repo.EXPECT().Delete(gomock.Any(), nodePath).Return(nil)
	err = registry1.Deregister(&node)
	assert.Nil(t, err)

	err = registry1.Close()
	assert.NoError(t, err)

	registry1 = NewRegistry(repo, constants.LiveNodesPath, 100)
	err = registry1.Close()
	assert.NoError(t, err)

	r := registry1.(*registry)
	r.register("/data/pant", &node)

	registry1 = NewRegistry(repo, constants.LiveNodesPath, 100)
	r = registry1.(*registry)

	// cancel ctx in timer
	time.AfterFunc(100*time.Millisecond, func() {
		r.cancel()
	})
	r.register("/data/pant", &node)
}
