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

package rpc

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
)

func TestConnectionManager_CreateConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskClientFct := NewMockTaskClientFactory(ctrl)
	connection := NewConnectionManager(taskClientFct)
	target := &models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}

	testCases := []struct {
		desc    string
		target  models.Node
		prepare func()
	}{
		{
			desc:   "create task client failure",
			target: &models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
			prepare: func() {
				taskClientFct.EXPECT().CreateTaskClient(gomock.Any()).Return(fmt.Errorf("err"))
			},
		},
		{
			desc:   "create task client successfully",
			target: &models.StatelessNode{HostIP: "1.1.1.2", GRPCPort: 9000},
			prepare: func() {
				taskClientFct.EXPECT().CreateTaskClient(gomock.Any()).Return(nil)
			},
		},
		{
			desc:   "connection exist",
			target: target,
			prepare: func() {
				taskClientFct.EXPECT().CreateTaskClient(gomock.Any()).Return(nil)
				connection.CreateConnection(target)
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(_ *testing.T) {
			tt.prepare()

			connection.CreateConnection(tt.target)
		})
	}
}

func TestConnectionManager_CloseConnection(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskClientFct := NewMockTaskClientFactory(ctrl)
	connection := NewConnectionManager(taskClientFct)
	target := &models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}

	testCases := []struct {
		desc    string
		target  models.Node
		prepare func()
	}{
		{
			desc:   "close non-existent connection",
			target: &models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000},
			prepare: func() {
				taskClientFct.EXPECT().CloseTaskClient(gomock.Any()).Return(false, nil)
			},
		},
		{
			desc:   "close task client successfully",
			target: &models.StatelessNode{HostIP: "1.1.1.2", GRPCPort: 9000},
			prepare: func() {
				taskClientFct.EXPECT().CloseTaskClient(gomock.Any()).Return(true, nil)
			},
		},
		{
			desc:   "close task client failure",
			target: target,
			prepare: func() {
				taskClientFct.EXPECT().CloseTaskClient(gomock.Any()).Return(true, fmt.Errorf("err"))
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(_ *testing.T) {
			tt.prepare()

			connection.CloseConnection(tt.target)
		})
	}
}

func TestConnectionManager_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskClientFct := NewMockTaskClientFactory(ctrl)
	connection := NewConnectionManager(taskClientFct)
	target := &models.StatelessNode{HostIP: "1.1.1.1", GRPCPort: 9000}

	taskClientFct.EXPECT().CreateTaskClient(gomock.Any()).Return(nil)
	taskClientFct.EXPECT().CloseTaskClient(gomock.Any()).Return(true, nil)
	connection.CreateConnection(target)
	assert.NoError(t, connection.Close())
}
