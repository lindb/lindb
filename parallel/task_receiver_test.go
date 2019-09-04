package parallel

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/tsdb/series"
)

func TestTaskReceiver_Receive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobManager := NewMockJobManager(ctrl)
	taskManager := NewMockTaskManager(ctrl)
	jobManager.EXPECT().GetTaskManager().Return(taskManager).AnyTimes()

	receiver := NewTaskReceiver(jobManager)
	taskManager.EXPECT().Get("taskID").Return(nil)
	err := receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.Nil(t, err)

	taskManager.EXPECT().Complete("taskID")
	taskManager.EXPECT().Get("taskID").
		Return(newTaskContext("taskID", RootTask, "parentTaskID", "parentNode", 1))

	jobManager.EXPECT().GetJob(gomock.Any()).Return(NewJobContext(make(chan series.GroupedIterator), nil))
	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.Nil(t, err)

}
