package broker

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
)

func TestStorageClusterState_SetState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	taskClientFactory := rpc.NewMockTaskClientFactory(ctrl)
	taskClientFactory.EXPECT().CloseTaskClient(gomock.Any()).AnyTimes()

	state := newStorageClusterState(taskClientFactory)

	storageState := models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}})

	taskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(nil)
	state.SetState(storageState)
	assert.Equal(t, 1, len(state.taskStreams))
	assert.NotNil(t, state.taskStreams["1.1.1.1:9000"])

	storageState = models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.2", Port: 9000}})
	taskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(fmt.Errorf("err"))
	state.SetState(storageState)
	assert.Equal(t, 0, len(state.taskStreams))

	storageState = models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.2", Port: 9000}})
	taskClientFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(nil)
	state.SetState(storageState)
	assert.Equal(t, 1, len(state.taskStreams))

	state.close()
	assert.Equal(t, 0, len(state.taskStreams))
}
