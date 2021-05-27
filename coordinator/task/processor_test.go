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

	"github.com/lindb/lindb/pkg/state"
)

func TestTaskProcessor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	txn := state.NewMockTransaction(ctrl)
	txn.EXPECT().ModRevisionCmp(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	txn.EXPECT().Put(gomock.Any(), gomock.Any()).AnyTimes()
	repo.EXPECT().NewTransaction().Return(txn).AnyTimes()

	proc := NewMockProcessor(ctrl)
	proc.EXPECT().Kind().Return(Kind("test")).AnyTimes()
	proc.EXPECT().Concurrency().Return(0)
	proc.EXPECT().RetryBackOff().Return(time.Duration(0))
	proc.EXPECT().RetryCount().Return(0)

	// submit fail
	taskProc := newTaskProcessor(context.TODO(), proc, repo)
	err := taskProc.Submit(taskEvent{task: Task{Kind: "tt"}})
	assert.NotNil(t, err)

	// submit task
	proc.EXPECT().Kind().Return(Kind("test")).AnyTimes()
	proc.EXPECT().Process(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	repo.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))

	err = taskProc.Submit(taskEvent{task: Task{Kind: "test"}})
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(100 * time.Millisecond)
	taskProc.Stop()

	err = taskProc.Submit(taskEvent{task: Task{Kind: "test"}})
	assert.NotNil(t, err)
}

func TestTaskProcessor_process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := state.NewMockRepository(ctrl)
	txn := state.NewMockTransaction(ctrl)
	txn.EXPECT().ModRevisionCmp(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	txn.EXPECT().Put(gomock.Any(), gomock.Any()).AnyTimes()
	repo.EXPECT().NewTransaction().Return(txn).AnyTimes()

	proc := NewMockProcessor(ctrl)
	proc.EXPECT().Kind().Return(Kind("test")).AnyTimes()
	proc.EXPECT().Concurrency().Return(0)
	proc.EXPECT().RetryBackOff().Return(time.Duration(0))
	proc.EXPECT().RetryCount().Return(0)

	taskProc := newTaskProcessor(context.TODO(), proc, repo)
	proc.EXPECT().Process(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	repo.EXPECT().Commit(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	taskProc.wg.Add(1)
	taskProc.process(taskEvent{task: Task{Kind: "test"}})
}
