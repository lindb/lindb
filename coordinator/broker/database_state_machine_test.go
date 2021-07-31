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

package broker

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
)

func TestNewDatabaseStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()

	// case 1: discovery err
	discovery1.EXPECT().Discovery(true).Return(fmt.Errorf("err"))
	_, err := NewDatabaseStateMachine(context.TODO(), factory)
	assert.Error(t, err)

	// case 2: normal case
	discovery1.EXPECT().Discovery(true).Return(nil)
	stateMachine, err := NewDatabaseStateMachine(context.TODO(), factory)
	assert.NoError(t, err)
	assert.NotNil(t, stateMachine)
}

func TestDBStateMachine_listen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	// normal case
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	discovery1.EXPECT().Discovery(true).Return(nil)
	stateMachine, err := NewDatabaseStateMachine(context.TODO(), factory)
	assert.NoError(t, err)

	db := models.Database{Name: "test"}
	data := encoding.JSONMarshal(&db)
	stateMachine.OnCreate("/data/test", data)

	db = models.Database{Name: "test3"}
	data = encoding.JSONMarshal(&db)
	stateMachine.OnCreate("/data/test3", data)

	db2, ok := stateMachine.GetDatabaseCfg("test")
	assert.True(t, ok)
	assert.Equal(t, models.Database{Name: "test"}, db2)

	stateMachine.OnCreate("/data/test2", []byte{1, 1})
	_, ok = stateMachine.GetDatabaseCfg("test2")
	assert.False(t, ok)

	data = encoding.JSONMarshal(&models.Database{})
	stateMachine.OnCreate("/data/test2", data)
	_, ok = stateMachine.GetDatabaseCfg("test2")
	assert.False(t, ok)

	stateMachine.OnDelete("/data/test")
	_, ok = stateMachine.GetDatabaseCfg("test")
	assert.False(t, ok)

	db2, ok = stateMachine.GetDatabaseCfg("test3")
	assert.True(t, ok)
	assert.Equal(t, db, db2)

	discovery1.EXPECT().Close()
	_ = stateMachine.Close()

	_, ok = stateMachine.GetDatabaseCfg("test3")
	assert.False(t, ok)
}
