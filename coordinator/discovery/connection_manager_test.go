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
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
)

func Test_ConnectionManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTaskClientFactory := rpc.NewMockTaskClientFactory(ctrl)
	cm := &ConnectionManager{
		RoleFrom:          "broker",
		RoleTo:            "broker",
		Connections:       make(map[string]struct{}),
		TaskClientFactory: mockTaskClientFactory,
	}
	mockTaskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(nil).Times(3)
	cm.CreateConnection(models.Node{IP: "192.168.1.1", Port: 1000})
	cm.CreateConnection(models.Node{IP: "192.168.1.2", Port: 2000})
	cm.CreateConnection(models.Node{IP: "192.168.1.3", Port: 3000})
	assert.Len(t, cm.Connections, 3)
	mockTaskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(io.ErrClosedPipe).Times(2)
	cm.CreateConnection(models.Node{IP: "192.168.1.3", Port: 4000})
	cm.CreateConnection(models.Node{IP: "192.168.1.3", Port: 3000})
	assert.Len(t, cm.Connections, 2)

	mockTaskClientFactory.EXPECT().CloseTaskClient(gomock.Any()).
		Return(true, nil).AnyTimes()
	cm.CloseInactiveNodeConnections([]string{
		"192.168.1.1:9000",
		"192.168.1.1:1000",
		"192.168.1.2:2000",
	})
	assert.Len(t, cm.Connections, 2)

	cm.CloseInactiveNodeConnections([]string{"192.168.1.1:1000"})
	assert.Len(t, cm.Connections, 1)

	cm.CloseInactiveNodeConnections([]string{})
	assert.Len(t, cm.Connections, 0)

}
