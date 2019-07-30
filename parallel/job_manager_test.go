package parallel

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestJobManager_SubmitJob(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskHandler := pb.NewMockTaskService_HandleClient(ctrl)
	taskManager := NewMockTaskManager(ctrl)
	taskManager.EXPECT().Submit(gomock.Any()).AnyTimes()
	taskManager.EXPECT().AllocTaskID().Return("TaskID").AnyTimes()
	taskSender := NewMockTaskSenderManager(ctrl)
	taskManager.EXPECT().GetTaskSenderManager().Return(nil)

	jobManager := NewJobManager(taskManager)
	physicalPlan := models.NewPhysicalPlan(models.Root{Indicator: "1.1.1.3:8000", NumOfTask: 1})
	physicalPlan.AddLeaf(models.Leaf{
		BaseNode: models.BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.1:9000",
		},
		ShardIDs: []int32{1, 2, 4},
	})
	err := jobManager.SubmitJob(physicalPlan)
	assert.Equal(t, errNoTaskSender, err)

	taskSender.EXPECT().GetClientStream(gomock.Any()).Return(nil)
	taskManager.EXPECT().GetTaskSenderManager().Return(taskSender)
	err = jobManager.SubmitJob(physicalPlan)
	assert.Equal(t, errNoSendStream, err)

	taskManager.EXPECT().GetTaskSenderManager().Return(taskSender).AnyTimes()
	taskSender.EXPECT().GetClientStream(gomock.Any()).Return(taskHandler).AnyTimes()
	taskHandler.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("error"))
	err = jobManager.SubmitJob(physicalPlan)
	assert.Equal(t, errTaskSend, err)

	taskHandler.EXPECT().Send(gomock.Any()).AnyTimes()
	err = jobManager.SubmitJob(physicalPlan)
	assert.Nil(t, err)
}

func TestJobManager_SubmitJob_2(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskHandler := pb.NewMockTaskService_HandleClient(ctrl)
	taskManager := NewMockTaskManager(ctrl)
	taskManager.EXPECT().Submit(gomock.Any()).AnyTimes()
	taskManager.EXPECT().AllocTaskID().Return("TaskID").AnyTimes()
	taskSender := NewMockTaskSenderManager(ctrl)
	taskManager.EXPECT().GetTaskSenderManager().Return(nil)

	jobManager := NewJobManager(taskManager)
	physicalPlan := models.NewPhysicalPlan(models.Root{Indicator: "1.1.1.3:8000", NumOfTask: 1})
	physicalPlan.AddIntermediate(models.Intermediate{
		BaseNode: models.BaseNode{
			Parent:    "1.1.1.3:8000",
			Indicator: "1.1.1.1:9000",
		},
	})
	err := jobManager.SubmitJob(physicalPlan)
	assert.Equal(t, errNoTaskSender, err)

	taskSender.EXPECT().GetClientStream(gomock.Any()).Return(nil)
	taskManager.EXPECT().GetTaskSenderManager().Return(taskSender)
	err = jobManager.SubmitJob(physicalPlan)
	assert.Equal(t, errNoSendStream, err)

	taskManager.EXPECT().GetTaskSenderManager().Return(taskSender).AnyTimes()
	taskSender.EXPECT().GetClientStream(gomock.Any()).Return(taskHandler).AnyTimes()
	taskHandler.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("error"))
	err = jobManager.SubmitJob(physicalPlan)
	assert.Equal(t, errTaskSend, err)

	taskHandler.EXPECT().Send(gomock.Any()).AnyTimes()
	err = jobManager.SubmitJob(physicalPlan)
	assert.Nil(t, err)
}
