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

package elect

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/lindb/common/pkg/encoding"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
)

func TestElection_Initialize(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	eventCh := make(chan *state.Event)

	repo := state.NewMockRepository(ctrl)

	listener1 := NewMockListener(ctrl)
	node1 := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(nil)
	election := NewElection(context.TODO(), repo, &node1, 1, listener1)
	election.Initialize()
	election.Close()

	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(eventCh)
	election = NewElection(context.TODO(), repo, &node1, 1, listener1)
	election.Initialize()
	election.Close()
}

func TestElection_Elect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	listener1 := NewMockListener(ctrl)

	node1 := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(nil)
	election := NewElection(context.TODO(), repo, &node1, 1, listener1)
	election.Initialize()
	election.Elect()
	election.Close()
}

func TestElection_elect(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	repo := state.NewMockRepository(ctrl)
	listener1 := NewMockListener(ctrl)

	node1 := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(nil)
	election1 := NewElection(ctx, repo, &node1, 1, listener1)
	election1.Initialize()
	e := election1.(*election)
	time.AfterFunc(700*time.Millisecond, func() {
		e.retryCh <- 1
		time.Sleep(10 * time.Millisecond)
		close(e.retryCh)
		cancel()
	})

	// fail
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(false, nil, fmt.Errorf("err"))
	// success
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	e.elect()

	election1.Close()
}

func TestElection_Is_Follower(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx, cancel := context.WithCancel(context.TODO())
	repo := state.NewMockRepository(ctrl)
	listener1 := NewMockListener(ctrl)

	node1 := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(nil)
	election1 := NewElection(ctx, repo, &node1, 1, listener1)
	election1.Initialize()
	e := election1.(*election)
	time.AfterFunc(700*time.Millisecond, func() {
		e.retryCh <- 1
		time.Sleep(10 * time.Millisecond)
		close(e.retryCh)
		cancel()
	})

	// success
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(false, nil, nil).AnyTimes()
	e.elect()

	election1.Close()
}

func TestElection_Handle_Event(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)

	listener1 := NewMockListener(ctrl)
	node1 := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(nil)
	election1 := NewElection(context.TODO(), repo, &node1, 1, listener1)
	election1.Initialize()
	e := election1.(*election)

	eventCh := make(chan *state.Event)

	go func() {
		sendEvent(eventCh, &state.Event{
			Type: state.EventTypeModify,
			KeyValues: []state.EventKeyValue{
				{Key: constants.MasterPath, Value: []byte{1, 1, 2}},
			},
		})
		sendEvent(eventCh, &state.Event{
			Type: state.EventTypeAll,
		})
		sendEvent(eventCh, &state.Event{
			Type: state.EventTypeModify,
			Err:  fmt.Errorf("err"),
		})
		data := encoding.JSONMarshal(&models.Master{Node: &node1})
		listener1.EXPECT().OnFailOver().Return(nil)
		sendEvent(eventCh, &state.Event{
			Type: state.EventTypeModify,
			KeyValues: []state.EventKeyValue{
				{Key: constants.MasterPath, Value: data},
			},
		})

		assert.True(t, e.IsMaster())

		// close chan
		close(eventCh)
	}()

	e.handleMasterChange(eventCh)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(nil)
	election1.Close()
}

func TestElection_handle_event(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	listener1 := NewMockListener(ctrl)

	node1 := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 2080}
	repo.EXPECT().Watch(gomock.Any(), gomock.Any(), true).Return(nil)
	election1 := NewElection(context.TODO(), repo, &node1, 1, listener1)
	assert.Nil(t, election1.GetMaster())
	election1.Initialize()
	e := election1.(*election)
	data := encoding.JSONMarshal(&models.Master{Node: &node1})

	time.AfterFunc(10*time.Millisecond, func() {
		<-e.retryCh
	})
	listener1.EXPECT().OnFailOver().Return(fmt.Errorf("err"))
	e.handleEvent(&state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.False(t, e.IsMaster())
	assert.Nil(t, e.GetMaster())

	listener1.EXPECT().OnFailOver().Return(nil)
	e.handleEvent(&state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: constants.MasterPath, Value: data},
		},
	})
	assert.True(t, e.IsMaster())
	assert.Equal(t, &node1, e.GetMaster().Node)

	time.AfterFunc(100*time.Millisecond, func() {
		<-e.retryCh
	})

	listener1.EXPECT().OnResignation()
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	e.handleEvent(&state.Event{
		Type: state.EventTypeDelete,
	})
	assert.False(t, e.IsMaster())
	assert.Equal(t, models.Master{}, *e.GetMaster())

	election1.Close()
}

func TestElection_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	listener1 := NewMockListener(ctrl)

	node1 := models.StatelessNode{HostIP: "127.0.0.1", GRPCPort: 2080}
	election1 := NewElection(context.TODO(), repo, &node1, 1, listener1)
	election1.Close()
	e := election1.(*election)
	e.elect()
	election1.Close()

	election1 = NewElection(context.TODO(), repo, &node1, 1, listener1)

	time.AfterFunc(100*time.Millisecond, func() {
		election1.Close()
	})
	repo.EXPECT().Elect(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(true, nil, nil).AnyTimes()
	e = election1.(*election)
	e.elect()
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(100 * time.Millisecond)
}
