package parallel

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	pb "github.com/lindb/lindb/rpc/proto/common"
	"github.com/lindb/lindb/service"
)

func TestLeafProcessor_Process(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storageService := service.NewMockStorageService(ctrl)
	executorFactory := NewMockExecutorFactory(ctrl)

	currentNode := models.Node{IP: "1.1.1.3", Port: 8000}
	processor := NewLeafTask(currentNode, storageService, executorFactory)
	// unmarshal error
	err := processor.Process(&pb.TaskRequest{PhysicalPlan: nil})
	assert.Equal(t, errUnmarshalPlan, err)

	plan, _ := json.Marshal(&models.PhysicalPlan{
		Leafs: []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.4:8000"}}},
	})
	// wrong request
	err = processor.Process(&pb.TaskRequest{PhysicalPlan: plan})
	assert.Equal(t, errWrongRequest, err)

	plan, _ = json.Marshal(&models.PhysicalPlan{
		Database: "test_db",
		Leafs:    []models.Leaf{{BaseNode: models.BaseNode{Indicator: "1.1.1.3:8000"}}},
	})
	// db not exist
	storageService.EXPECT().GetEngine(gomock.Any()).Return(nil)
	err = processor.Process(&pb.TaskRequest{PhysicalPlan: plan})
	assert.Equal(t, errNoDatabase, err)
}
