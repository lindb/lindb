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

package task

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
)

func TestController_Submit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := NewControllerFactory()
	repo := state.NewMockRepository(ctrl)
	txn := state.NewMockTransaction(ctrl)

	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), true).Return(nil)
	controller := factory.CreateController(context.TODO(), repo)
	assert.NotNil(t, controller)
	params := make([]ControllerTaskParam, 127)
	err := controller.Submit("", "", params)
	assert.Equal(t, ErrMaxTasksLimitExceeded, err)
	err = controller.Submit("", "", nil)
	assert.Nil(t, err)

	node1 := &models.Node{IP: "1.1.1.1", Port: 8000}
	node2 := &models.Node{IP: "1.1.1.1", Port: 8000}
	repo.EXPECT().NewTransaction().Return(txn)
	txn.EXPECT().Put(gomock.Any(), gomock.Any()).MaxTimes(3)
	repo.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(nil)
	err = controller.Submit(kindDummy, "wtf-1019-07-05--1", []ControllerTaskParam{
		{NodeID: node1.Indicator(), Params: dummyParams{}},
		{NodeID: node2.Indicator(), Params: dummyParams{}},
	})
	if err != nil {
		t.Fatal(err)
	}

	repo.EXPECT().NewTransaction().Return(txn)
	txn.EXPECT().Put(gomock.Any(), gomock.Any()).MaxTimes(3)
	repo.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = controller.Submit(kindDummy, "wtf-1019-07-05--1", []ControllerTaskParam{
		{NodeID: node1.Indicator(), Params: dummyParams{}},
		{NodeID: node2.Indicator(), Params: dummyParams{}},
	})
	assert.NotNil(t, err)

	assert.Equal(t, taskCoordinatorKey+"/executor/node/kinds/k/names/name",
		controller.taskKey("k", "name", "node"))
	assert.Equal(t, taskCoordinatorKey+"/status/kinds/k/names/name",
		controller.statusKey("k", "name"))

	// wait goroutine exit
	time.Sleep(100 * time.Millisecond)
	err = controller.Close()
	if err != nil {
		t.Fatal(err)
	}

	err = controller.Submit("", "", []ControllerTaskParam{})
	assert.Equal(t, ErrControllerClosed, err)
}

func TestController_run(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := NewControllerFactory()
	repo := state.NewMockRepository(ctrl)

	// test event channel close
	eventCh := make(chan *state.Event)
	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), true).Return(eventCh)
	controller := factory.CreateController(context.TODO(), repo)
	assert.NotNil(t, controller)
	close(eventCh)
	// wait goroutine exit
	time.Sleep(100 * time.Millisecond)
	err := controller.Close()
	if err != nil {
		t.Fatal(err)
	}

	// normal case
	eventCh = make(chan *state.Event)
	repo.EXPECT().WatchPrefix(gomock.Any(), gomock.Any(), true).Return(eventCh)
	controller = factory.CreateController(context.TODO(), repo)
	assert.NotNil(t, controller)

	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeAll,
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeDelete,
	})
	sendEvent(eventCh, &state.Event{
		Err: fmt.Errorf("err"),
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.taskKey("k", "name", "node"), Value: []byte{1, 1, 1}},
		},
	})
	task := Task{}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.taskKey("k", "name", "node"), Value: encoding.JSONMarshal(&task)},
		},
	})
	task = Task{State: StateDoneOK}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.taskKey("k", "name", "node"), Value: encoding.JSONMarshal(&task)},
		},
	})
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.statusKey("k", "name"), Value: []byte{1, 13}},
		},
	})
	taskGroup := groupedTasks{
		State: StateDoneOK,
		Tasks: []Task{{Kind: "test"}},
	}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.statusKey("k", "name"), Value: encoding.JSONMarshal(&taskGroup)},
		},
	})
	taskGroup = groupedTasks{
		Tasks: []Task{{Kind: "test", Name: "test-name", Executor: "node"}},
	}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.statusKey("test", "name"), Value: encoding.JSONMarshal(&taskGroup)},
		},
	})
	task = Task{Kind: "test", State: StateDoneOK, Name: "no-name"}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.taskKey("test", "no-name", "node"), Value: encoding.JSONMarshal(&task)},
		},
	})
	// update task error
	task = Task{Kind: "test", State: StateDoneErr, Executor: "node", Name: "test-name", ErrMsg: "err-msg"}
	txn := state.NewMockTransaction(ctrl)
	repo.EXPECT().NewTransaction().Return(txn).MaxTimes(2)
	txn.EXPECT().Delete(gomock.Any()).MaxTimes(2)
	txn.EXPECT().Put(gomock.Any(), gomock.Any()).MaxTimes(2)
	txn.EXPECT().ModRevisionCmp(gomock.Any(), gomock.Any(), gomock.Any()).MaxTimes(2)
	repo.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.taskKey("test", "test-name", "node"), Value: encoding.JSONMarshal(&task)},
		},
	})
	// update status success
	taskGroup = groupedTasks{
		Tasks: []Task{{Kind: "test-kind", Name: "test-name", Executor: "node"}},
	}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.statusKey("test-kind", "name"), Value: encoding.JSONMarshal(&taskGroup)},
		},
	})
	task = Task{Kind: "test-kind", State: StateDoneErr, Name: "test-name", Executor: "node", ErrMsg: "err-msg"}
	repo.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(nil)
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.taskKey("test-kind", "test-name", "node"), Value: encoding.JSONMarshal(&task)},
		},
	})

	// no executor
	taskGroup = groupedTasks{
		Tasks: []Task{{Kind: "test-kind", Name: "test-name"}},
	}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.statusKey("test-kind", "name"), Value: encoding.JSONMarshal(&taskGroup)},
		},
	})
	task = Task{Kind: "test-kind", State: StateDoneErr, Name: "test-name", Executor: "node", ErrMsg: "err-msg"}
	sendEvent(eventCh, &state.Event{
		Type: state.EventTypeModify,
		KeyValues: []state.EventKeyValue{
			{Key: controller.taskKey("test-kind", "test-name", "node"), Value: encoding.JSONMarshal(&task)},
		},
	})

	// wait goroutine exit
	time.Sleep(100 * time.Millisecond)
	err = controller.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func sendEvent(eventCh chan *state.Event, event *state.Event) {
	eventCh <- event
	time.Sleep(10 * time.Millisecond)
}
