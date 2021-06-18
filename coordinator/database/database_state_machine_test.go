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

package database

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestNewDBStateMachine(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	_, err := NewDBStateMachine(context.TODO(), repo, factory)
	assert.Error(t, err)

	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err = NewDBStateMachine(context.TODO(), repo, factory)
	assert.Error(t, err)

	// normal case
	data, _ := json.Marshal(&models.Database{Name: "test"})
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]state.KeyValue{
		{Value: data},
		{Value: []byte{1, 1, 2}},
	}, nil)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(nil)
	stateMachine, err := NewDBStateMachine(context.TODO(), repo, factory)
	assert.NoError(t, err)
	assert.NotNil(t, stateMachine)
}

func TestDBStateMachine_listen(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	factory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	// normal case
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	factory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1)
	discovery1.EXPECT().Discovery().Return(nil)
	stateMachine, err := NewDBStateMachine(context.TODO(), repo, factory)
	assert.NoError(t, err)

	db := models.Database{Name: "test"}
	data, _ := json.Marshal(&db)
	stateMachine.OnCreate("/data/test", data)

	db2, ok := stateMachine.GetDatabaseCfg("test")
	assert.True(t, ok)
	assert.Equal(t, db, db2)

	stateMachine.OnCreate("/data/test2", []byte{1, 1})
	_, ok = stateMachine.GetDatabaseCfg("test2")
	assert.False(t, ok)

	data, _ = json.Marshal(&models.Database{})
	stateMachine.OnCreate("/data/test2", data)
	_, ok = stateMachine.GetDatabaseCfg("test2")
	assert.False(t, ok)

	stateMachine.OnDelete("/data/test")
	_, ok = stateMachine.GetDatabaseCfg("test")
	assert.False(t, ok)

	discovery1.EXPECT().Close()
	_ = stateMachine.Close()
}
