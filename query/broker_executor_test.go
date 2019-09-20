package query

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/parallel"
)

func TestBrokerExecutor_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	currentNode := generateBrokerActiveNode("1.1.1.3", 8000)

	nodeStateMachine := broker.NewMockNodeStateMachine(ctrl)
	nodeStateMachine.EXPECT().GetCurrentNode().Return(currentNode.Node).AnyTimes()
	replicaStateMachine := replica.NewMockStatusStateMachine(ctrl)
	jobManager := parallel.NewMockJobManager(ctrl)

	exec := newBrokerExecutor(context.TODO(), "test_db", "select f from cpu",
		replicaStateMachine, nodeStateMachine, jobManager)
	replicaStateMachine.EXPECT().GetQueryableReplicas("test_db").Return(nil)
	_ = exec.Execute()
	assert.Equal(t, errNoAvailableStorageNode, exec.Error())

	storageNodes := map[string][]int32{
		"1.1.1.1:9000": {1, 2, 4},
		"1.1.1.2:9000": {3, 6, 9},
		"1.1.1.3:9000": {5, 7, 8},
		"1.1.1.4:9000": {10, 13, 15},
		"1.1.1.5:9000": {11, 12, 14},
	}
	brokerNodes := []models.ActiveNode{
		generateBrokerActiveNode("1.1.1.1", 8000),
		generateBrokerActiveNode("1.1.1.2", 8000),
		currentNode,
		generateBrokerActiveNode("1.1.1.4", 8000),
	}
	exec = newBrokerExecutor(context.TODO(), "test_db", "select f fro",
		replicaStateMachine, nodeStateMachine, jobManager)
	replicaStateMachine.EXPECT().GetQueryableReplicas("test_db").Return(storageNodes)
	nodeStateMachine.EXPECT().GetActiveNodes().Return(brokerNodes)
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())

	exec = newBrokerExecutor(context.TODO(), "test_db", "select f from cpu",
		replicaStateMachine, nodeStateMachine, jobManager)
	replicaStateMachine.EXPECT().GetQueryableReplicas("test_db").Return(storageNodes)
	nodeStateMachine.EXPECT().GetActiveNodes().Return(brokerNodes)
	jobManager.EXPECT().SubmitJob(gomock.Any())
	_ = exec.Execute()
	assert.Nil(t, exec.Error())
	assert.NotNil(t, exec.Statement())

	// submit job error
	exec = newBrokerExecutor(context.TODO(), "test_db", "select f from cpu",
		replicaStateMachine, nodeStateMachine, jobManager)
	replicaStateMachine.EXPECT().GetQueryableReplicas("test_db").Return(storageNodes)
	nodeStateMachine.EXPECT().GetActiveNodes().Return(brokerNodes)
	jobManager.EXPECT().SubmitJob(gomock.Any()).Return(errors.New("submit job error"))
	_ = exec.Execute()
	assert.NotNil(t, exec.Error())
}
