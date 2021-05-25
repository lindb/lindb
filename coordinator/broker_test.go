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

package coordinator

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/database"
	"github.com/lindb/lindb/coordinator/replica"
)

func TestBrokerStateMachines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := NewMockStateMachineFactory(ctrl)
	brokerSMs := NewBrokerStateMachines(factory)
	nodeSM := broker.NewMockNodeStateMachine(ctrl)
	replicaSM := replica.NewMockStatusStateMachine(ctrl)
	storageStateSM := broker.NewMockStorageStateMachine(ctrl)
	replicatorSM := replica.NewMockReplicatorStateMachine(ctrl)
	dbSM := database.NewMockDBStateMachine(ctrl)

	factory.EXPECT().CreateNodeStateMachine().Return(nil, fmt.Errorf("err"))
	err := brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateNodeStateMachine().Return(nodeSM, nil).AnyTimes()
	factory.EXPECT().CreateReplicatorStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateReplicatorStateMachine().Return(replicatorSM, nil).AnyTimes()
	factory.EXPECT().CreateStorageStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateStorageStateMachine().Return(storageStateSM, nil).AnyTimes()
	factory.EXPECT().CreateReplicaStatusStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateReplicaStatusStateMachine().Return(replicaSM, nil).AnyTimes()
	factory.EXPECT().CreateDatabaseStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateDatabaseStateMachine().Return(dbSM, nil).AnyTimes()
	err = brokerSMs.Start()
	assert.NoError(t, err)

	nodeSM.EXPECT().Close().Return(fmt.Errorf("err"))
	replicaSM.EXPECT().Close().Return(fmt.Errorf("err"))
	storageStateSM.EXPECT().Close().Return(fmt.Errorf("err"))
	replicatorSM.EXPECT().Close().Return(fmt.Errorf("err"))
	dbSM.EXPECT().Close().Return(fmt.Errorf("err"))
	brokerSMs.Stop()
}
