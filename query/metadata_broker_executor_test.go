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

package query

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/sql/stmt"
)

func TestMetadataBrokerExecutor_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	nodeStateMachine := discovery.NewMockActiveNodeStateMachine(ctrl)
	replicaStateMachine := broker.NewMockReplicaStatusStateMachine(ctrl)
	jobManager := parallel.NewMockJobManager(ctrl)

	exec := newMetadataBrokerExecutor(context.TODO(), "test_db", &stmt.Metadata{},
		nodeStateMachine, replicaStateMachine, jobManager)

	// no storage node
	replicaStateMachine.EXPECT().GetQueryableReplicas(gomock.Any()).Return(nil)
	rs, err := exec.Execute()
	assert.Error(t, err)
	assert.Nil(t, rs)

	// submit job err
	nodeStateMachine.EXPECT().GetCurrentNode().Return(models.Node{IP: "2.2.2.2", Port: 1234})
	replicaStateMachine.EXPECT().GetQueryableReplicas(gomock.Any()).Return(map[string][]int32{"1.1.1.1:1234": {1, 2, 3}})
	jobManager.EXPECT().SubmitMetadataJob(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	rs, err = exec.Execute()
	assert.Error(t, err)
	assert.Nil(t, rs)

	// execute query job
	jobManager.EXPECT().SubmitMetadataJob(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	e := exec.(*metadataBrokerExecutor)
	resultCh := make(chan []string)
	go func() {
		resultCh <- []string{"b", "d", "a"}
		close(resultCh)
	}()
	rs, err = e.submitJob(nil, resultCh)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "d"}, rs)
}
