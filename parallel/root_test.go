package parallel

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestRootTask_Receive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskManager := NewMockTaskManager(ctrl)

	receiver := newRookTask(taskManager)
	taskManager.EXPECT().Get("taskID").Return(nil)
	err := receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.Nil(t, err)

	taskManager.EXPECT().Complete("taskID")
	taskManager.EXPECT().Get("taskID").
		Return(newTaskContext("taskID", "parentTaskID", "parentNode", 1))
	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.Nil(t, err)

}
