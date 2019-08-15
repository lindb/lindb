package coordinator

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/replica"
)

func TestBrokerStateMachines(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	factory := NewMockStateMachineFactory(ctrl)
	brokerSMs := NewBrokerStateMachines(factory)
	nodeSM := broker.NewMockNodeStateMachine(ctrl)
	replicaSM := replica.NewMockStatusStateMachine(ctrl)
	storageStateSM := broker.NewMockStorageStateMachine(ctrl)
	replicatorSM := replica.NewMockReplicatorStateMachine(ctrl)

	factory.EXPECT().CreateNodeStateMachine().Return(nil, fmt.Errorf("err"))
	err := brokerSMs.Start()
	assert.NotNil(t, err)

	factory.EXPECT().CreateNodeStateMachine().Return(nodeSM, nil)
	factory.EXPECT().CreateReplicatorStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.NotNil(t, err)

	factory.EXPECT().CreateNodeStateMachine().Return(nodeSM, nil)
	factory.EXPECT().CreateReplicatorStateMachine().Return(replicatorSM, nil)
	factory.EXPECT().CreateStorageStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.NotNil(t, err)

	factory.EXPECT().CreateNodeStateMachine().Return(nodeSM, nil)
	factory.EXPECT().CreateReplicatorStateMachine().Return(replicatorSM, nil)
	factory.EXPECT().CreateStorageStateMachine().Return(storageStateSM, nil)
	factory.EXPECT().CreateReplicaStatusStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.NotNil(t, err)

	factory.EXPECT().CreateNodeStateMachine().Return(nodeSM, nil)
	factory.EXPECT().CreateStorageStateMachine().Return(storageStateSM, nil)
	factory.EXPECT().CreateReplicaStatusStateMachine().Return(replicaSM, nil)
	factory.EXPECT().CreateReplicatorStateMachine().Return(replicatorSM, nil)
	err = brokerSMs.Start()
	if err != nil {
		t.Fatal(err)
	}

	nodeSM.EXPECT().Close().Return(fmt.Errorf("err"))
	replicaSM.EXPECT().Close().Return(fmt.Errorf("err"))
	storageStateSM.EXPECT().Close().Return(fmt.Errorf("err"))
	replicatorSM.EXPECT().Close().Return(fmt.Errorf("err"))
	brokerSMs.Stop()
}
