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

package command

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	depspkg "github.com/lindb/lindb/app/broker/deps"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/internal/client"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/sql/stmt"
)

func TestRequst(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		requestCli = client.NewRequestCli()
		ctrl.Finish()
	}()

	cli := client.NewMockRequestCli(ctrl)
	requestCli = cli

	stateMgr := broker.NewMockStateManager(ctrl)
	deps := &depspkg.HTTPDeps{
		StateMgr: stateMgr,
	}
	stateMgr.EXPECT().GetLiveNodes().Return([]models.StatelessNode{{
		HostIP:   "127.0.0.1",
		HTTPPort: 3000,
	}})
	cli.EXPECT().FetchRequestsByNodes(gomock.Any()).Return(nil)
	rs, err := RequestCommand(context.TODO(), deps, nil, &stmt.Request{})
	assert.NoError(t, err)
	assert.Nil(t, rs)
}
