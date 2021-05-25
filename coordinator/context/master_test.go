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

package context

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/storage"
)

func TestMasterContext_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	stateMachine := &StateMachine{}
	ctx := NewMasterContext(stateMachine)
	ctx.Close()

	admin := database.NewMockAdminStateMachine(ctrl)
	cluster := storage.NewMockClusterStateMachine(ctrl)

	stateMachine.DatabaseAdmin = admin
	stateMachine.StorageCluster = cluster

	admin.EXPECT().Close().Return(nil)
	cluster.EXPECT().Close().Return(nil)
	ctx.Close()

	admin.EXPECT().Close().Return(fmt.Errorf("err"))
	cluster.EXPECT().Close().Return(fmt.Errorf("err"))
	ctx.Close()
}
