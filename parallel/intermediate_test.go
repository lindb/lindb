package parallel

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestIntermediate_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskHandler := pb.NewMockTaskService_HandleClient(ctrl)
	taskManager := NewMockTaskManager(ctrl)
	taskManager.EXPECT().Submit(gomock.Any()).AnyTimes()
	taskSender := NewMockTaskSenderManager(ctrl)

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	processor := newIntermediateTask(currentNode, taskManager)

	// unmarshal error
	err := processor.Process(&pb.TaskRequest{PhysicalPlan: nil})
	assert.Equal(t, errUnmarshalPlan, err)

	plan, _ := json.Marshal(&models.PhysicalPlan{
		Intermediates: []models.Intermediate{{BaseNode: models.BaseNode{Indicator: "1.1.1.4:8000"}}},
	})
	// wrong request
	err = processor.Process(&pb.TaskRequest{PhysicalPlan: plan})
	assert.Equal(t, errWrongRequest, err)

	plan2, _ := json.Marshal(&models.PhysicalPlan{
		Intermediates: []models.Intermediate{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
		Leafs: []models.Leaf{
			{BaseNode: models.BaseNode{Parent: "1.1.1.3:8000", Indicator: "1.1.1.5:8000"}},
		},
	})
	// no send stream
	taskManager.EXPECT().AllocTaskID().Return("taskID").AnyTimes()
	taskManager.EXPECT().GetTaskSenderManager().Return(nil)
	err = processor.Process(&pb.TaskRequest{PhysicalPlan: plan2})
	assert.Equal(t, errNoTaskSender, err)

	taskManager.EXPECT().GetTaskSenderManager().Return(taskSender).AnyTimes()

	// no send stream
	taskSender.EXPECT().GetClientStream("1.1.1.5:8000").Return(nil)
	err = processor.Process(&pb.TaskRequest{PhysicalPlan: plan2})
	assert.Equal(t, errNoSendStream, err)

	// send error
	taskSender.EXPECT().GetClientStream("1.1.1.5:8000").Return(taskHandler)
	taskHandler.EXPECT().Send(gomock.Any()).Return(errTaskSend)
	err = processor.Process(&pb.TaskRequest{PhysicalPlan: plan2})
	assert.Equal(t, errTaskSend, err)

	// normal
	taskSender.EXPECT().GetClientStream("1.1.1.5:8000").Return(taskHandler)
	taskHandler.EXPECT().Send(gomock.Any()).Return(nil)
	err = processor.Process(&pb.TaskRequest{PhysicalPlan: plan2})
	if err != nil {
		t.Fatal(err)
	}

	// normal
	plan, _ = json.Marshal(&models.PhysicalPlan{
		Intermediates: []models.Intermediate{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	err = processor.Process(&pb.TaskRequest{PhysicalPlan: plan})
	if err != nil {
		t.Fatal(err)
	}
}

func TestIntermediate_Receive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskManager := NewMockTaskManager(ctrl)
	taskSenderManager := NewMockTaskSenderManager(ctrl)
	taskManager.EXPECT().GetTaskSenderManager().Return(taskSenderManager).AnyTimes()

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	receiver := newIntermediateTask(currentNode, taskManager)
	taskManager.EXPECT().Get("taskID").Return(nil)
	err := receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.Nil(t, err)

	// not found server stream
	taskManager.EXPECT().Complete("taskID")
	taskSenderManager.EXPECT().GetServerStream("parentNode").Return(nil)
	taskManager.EXPECT().Get("taskID").
		Return(newTaskContext("taskID", "parentTaskID", "parentNode", 1))
	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.Nil(t, err)

	// send task result error
	serverStream := pb.NewMockTaskService_HandleServer(ctrl)
	serverStream.EXPECT().Send(gomock.Any()).Return(fmt.Errorf("send error"))
	taskManager.EXPECT().Complete("taskID")
	taskSenderManager.EXPECT().GetServerStream("parentNode").Return(serverStream)
	taskManager.EXPECT().Get("taskID").
		Return(newTaskContext("taskID", "parentTaskID", "parentNode", 1))
	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.NotNil(t, err)

	// normal case
	serverStream.EXPECT().Send(gomock.Any()).Return(nil)
	taskManager.EXPECT().Complete("taskID")
	taskSenderManager.EXPECT().GetServerStream("parentNode").Return(serverStream)
	taskManager.EXPECT().Get("taskID").
		Return(newTaskContext("taskID", "parentTaskID", "parentNode", 1))
	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	assert.Nil(t, err)
}
