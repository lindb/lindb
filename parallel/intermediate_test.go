package parallel

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/sql"
)

func TestIntermediate_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskManager := NewMockTaskManager(ctrl)
	taskManager.EXPECT().Submit(gomock.Any()).AnyTimes()

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	processor := newIntermediateTask(currentNode, taskManager)

	// unmarshal error
	err := processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: nil})
	assert.Equal(t, errUnmarshalPlan, err)

	plan, _ := json.Marshal(&models.PhysicalPlan{
		Intermediates: []models.Intermediate{{BaseNode: models.BaseNode{Indicator: "1.1.1.4:8000"}}},
	})
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan})
	assert.Equal(t, errUnmarshalQuery, err)

	// wrong request
	query, _ := sql.Parse("select f from cpu where host='1.1.1.1' and time>'20190729 11:00:00' and time<'20190729 12:00:00'")
	data := encoding.JSONMarshal(query)
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.Equal(t, errWrongRequest, err)

	plan2, _ := json.Marshal(&models.PhysicalPlan{
		Intermediates: []models.Intermediate{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
		Leafs: []models.Leaf{
			{BaseNode: models.BaseNode{Parent: "1.1.1.3:8000", Indicator: "1.1.1.5:8000"}},
		},
	})
	taskManager.EXPECT().AllocTaskID().Return("taskID").AnyTimes()
	// send request error
	taskManager.EXPECT().SendRequest(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan2, Payload: data})
	assert.NotNil(t, err)

	// normal
	taskManager.EXPECT().SendRequest(gomock.Any(), gomock.Any()).Return(nil)
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan2, Payload: data})
	assert.NoError(t, err)

	// normal
	plan, _ = json.Marshal(&models.PhysicalPlan{
		Intermediates: []models.Intermediate{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	err = processor.Process(context.TODO(), &pb.TaskRequest{PhysicalPlan: plan, Payload: data})
	assert.NoError(t, err)
}

func TestIntermediate_Receive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskManager := NewMockTaskManager(ctrl)

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	receiver := newIntermediateTask(currentNode, taskManager)
	taskManager.EXPECT().Get("taskID").Return(nil)
	err := receiver.Receive(&pb.TaskResponse{TaskID: "taskID"})
	if err != nil {
		t.Fatal(err)
	}

	// send task result error
	merger := NewMockResultMerger(ctrl)
	merger.EXPECT().merge(gomock.Any())
	merger.EXPECT().close()
	taskManager.EXPECT().Complete("taskID")
	taskManager.EXPECT().Get("taskID").
		Return(newTaskContext("taskID", IntermediateTask, "parentTaskID", "parentNode", 1, merger))
	taskManager.EXPECT().SendResponse(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID", Completed: true})
	assert.NotNil(t, err)

	// normal case
	merger.EXPECT().merge(gomock.Any())
	merger.EXPECT().close()
	taskManager.EXPECT().Complete("taskID")
	taskManager.EXPECT().Get("taskID").
		Return(newTaskContext("taskID", IntermediateTask, "parentTaskID", "parentNode", 1, merger))
	taskManager.EXPECT().SendResponse(gomock.Any(), gomock.Any()).Return(nil)
	err = receiver.Receive(&pb.TaskResponse{TaskID: "taskID", Completed: true})
	if err != nil {
		t.Fatal(err)
	}
}
