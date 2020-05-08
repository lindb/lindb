package coordinator

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/database"
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
	dbSM := database.NewMockDBStateMachine(ctrl)

	factory.EXPECT().CreateNodeStateMachine().Return(nil, fmt.Errorf("err"))
	err := brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateNodeStateMachine().Return(nodeSM, nil).AnyTimes()
	factory.EXPECT().CreateReplicatorStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateReplicatorStateMachine().Return(replicatorSM, nil).AnyTimes()
	factory.EXPECT().CreateStorageStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateStorageStateMachine().Return(storageStateSM, nil).AnyTimes()
	factory.EXPECT().CreateReplicaStatusStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateReplicaStatusStateMachine().Return(replicaSM, nil).AnyTimes()
	factory.EXPECT().CreateDatabaseStateMachine().Return(nil, fmt.Errorf("err"))
	err = brokerSMs.Start()
	assert.Error(t, err)

	factory.EXPECT().CreateDatabaseStateMachine().Return(dbSM, nil).AnyTimes()
	err = brokerSMs.Start()
	assert.NoError(t, err)

	nodeSM.EXPECT().Close().Return(fmt.Errorf("err"))
	replicaSM.EXPECT().Close().Return(fmt.Errorf("err"))
	storageStateSM.EXPECT().Close().Return(fmt.Errorf("err"))
	replicatorSM.EXPECT().Close().Return(fmt.Errorf("err"))
	dbSM.EXPECT().Close().Return(fmt.Errorf("err"))
	brokerSMs.Stop()
}
