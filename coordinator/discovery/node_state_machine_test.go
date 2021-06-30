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
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
)

func TestNewNodeStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// case 1: discovery resource err
	factory := NewMockFactory(ctrl)
	discovery := NewMockDiscovery(ctrl)
	factory.EXPECT().CreateDiscovery(constants.NodesPath, gomock.Any()).Return(discovery).AnyTimes()
	discovery.EXPECT().Discovery(true).Return(fmt.Errorf("err"))
	sm, err := NewNodeStateMachine(context.TODO(), factory)
	assert.Error(t, err)
	assert.Nil(t, sm)

	// case 2: ok
	discovery.EXPECT().Discovery(true).Return(nil)
	sm, err = NewNodeStateMachine(context.TODO(), factory)
	assert.NoError(t, err)
	assert.NotNil(t, sm)

	// case 3: close
	discovery.EXPECT().Close()
	err = sm.Close()
	assert.NoError(t, err)
	_ = sm.Close()
	assert.Empty(t, sm.GetNodes())
}

func TestNodeStateMachine_Discovery(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// case 1: discovery resource err
	factory := NewMockFactory(ctrl)
	discovery := NewMockDiscovery(ctrl)
	factory.EXPECT().CreateDiscovery(constants.NodesPath, gomock.Any()).Return(discovery).AnyTimes()
	discovery.EXPECT().Discovery(true).Return(nil)
	sm, err := NewNodeStateMachine(context.TODO(), factory)
	assert.NoError(t, err)
	// case 1: unmarshal data err
	sm.OnCreate("/test", []byte{1, 2, 3})
	assert.Empty(t, sm.GetNodes())

	// case 2: discovery new node
	data, _ := json.Marshal(&models.Node{ID: 1, IP: "1.1.1.1", Port: 8080})
	sm.OnCreate("/1.1.1.1:8080", data)
	assert.Len(t, sm.GetNodes(), 1)

	// case 3: delete exist node
	sm.OnDelete("/1.1.1.1:8080")
	assert.Empty(t, sm.GetNodes())
}
