package broker

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/rpc"
	pb "github.com/lindb/lindb/rpc/proto/common"
)

func TestStorageClusterState_SetState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	streamFactory := rpc.NewMockClientStreamFactory(ctrl)
	clientStream := pb.NewMockTaskService_HandleClient(ctrl)
	clientStream.EXPECT().CloseSend().Return(fmt.Errorf("err")).AnyTimes()

	state := newStorageClusterState(streamFactory)

	storageState := models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.1", Port: 9000}})

	streamFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(clientStream, nil)
	state.SetState(storageState)
	assert.Equal(t, 1, len(state.taskStreams))
	assert.NotNil(t, state.taskStreams["1.1.1.1:9000"])

	storageState = models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.2", Port: 9000}})
	streamFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(nil, fmt.Errorf("err"))
	state.SetState(storageState)
	assert.Equal(t, 0, len(state.taskStreams))

	storageState = models.NewStorageState()
	storageState.AddActiveNode(&models.ActiveNode{Node: models.Node{IP: "1.1.1.2", Port: 9000}})
	streamFactory.EXPECT().CreateTaskClient(gomock.Any()).Return(clientStream, nil)
	state.SetState(storageState)
	assert.Equal(t, 1, len(state.taskStreams))

	state.close()
	assert.Equal(t, 0, len(state.taskStreams))
}
