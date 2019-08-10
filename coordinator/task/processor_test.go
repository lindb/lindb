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
	proc := NewMockProcessor(ctrl)
	proc.EXPECT().Kind().Return(Kind("test")).AnyTimes()
	proc.EXPECT().Concurrency().Return(0)
	proc.EXPECT().RetryBackOff().Return(time.Duration(0))
	proc.EXPECT().RetryCount().Return(0)

	taskProc := newTaskProcessor(context.TODO(), proc, repo)
	err := taskProc.Submit(taskEvent{task: Task{Kind: "tt"}})
	assert.NotNil(t, err)
	err = taskProc.Submit(taskEvent{task: Task{Kind: "test"}})
	if err != nil {
		t.Fatal(err)
	}

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
