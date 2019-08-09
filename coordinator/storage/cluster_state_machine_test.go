package storage

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/coordinator/task"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/service"
)

func TestClusterStateMachine_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	controllerFactory := task.NewMockControllerFactory(ctrl)
	storageService := service.NewMockStorageStateService(ctrl)
	shardAssignService := service.NewMockShardAssignService(ctrl)

	repoFactory := state.NewMockRepositoryFactory(ctrl)
	repo := state.NewMockRepository(ctrl)
	discoverFactory := discovery.NewMockFactory(ctrl)
	discovery1 := discovery.NewMockDiscovery(ctrl)
	cluster := NewMockCluster(ctrl)
	discoverFactory.EXPECT().CreateDiscovery(gomock.Any(), gomock.Any()).Return(discovery1).AnyTimes()
	clusterFactory := NewMockClusterFactory(ctrl)

	// list exist storage cluster err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, fmt.Errorf("err"))
	_, err := NewClusterStateMachine(context.TODO(), repo,
		controllerFactory, discoverFactory, clusterFactory, repoFactory,
		storageService, shardAssignService)

	assert.NotNil(t, err)

	// register discovery err
	repo.EXPECT().List(gomock.Any(), gomock.Any()).Return(nil, nil)
	discovery1.EXPECT().Discovery().Return(fmt.Errorf("err"))
	_, err = NewClusterStateMachine(context.TODO(), repo,
		controllerFactory, discoverFactory, clusterFactory, repoFactory,
		storageService, shardAssignService)
	assert.NotNil(t, err)

	// normal case
	repo.EXPECT().List(gomock.Any(), gomock.Any()).
		Return([][]byte{
			encoding.JSONMarshal(&models.StorageState{Name: "test1"}),
			{1, 2, 3},
			encoding.JSONMarshal(&models.StorageState{Name: "test2"}),
			encoding.JSONMarshal(&models.StorageState{Name: "test3"}),
			encoding.JSONMarshal(&models.StorageState{}),
		}, nil).AnyTimes()
	repo1 := state.NewMockRepository(ctrl)
	repo1.EXPECT().Close().Return(nil)

	gomock.InOrder(
		repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(repo1, nil),
		clusterFactory.EXPECT().newCluster(gomock.Any()).Return(nil, fmt.Errorf("err")),
		repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(state.NewMockRepository(ctrl), nil),
		clusterFactory.EXPECT().newCluster(gomock.Any()).Return(cluster, nil),
		repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(nil, fmt.Errorf("err")),
	)
	discovery1.EXPECT().Discovery().Return(nil)

	stateMachine, err := NewClusterStateMachine(context.TODO(), repo,
		controllerFactory, discoverFactory, clusterFactory, repoFactory,
		storageService, shardAssignService)

	assert.Nil(t, err)
	assert.NotNil(t, stateMachine)
	assert.Equal(t, 1, len(stateMachine.GetAllCluster()))
	assert.Equal(t, cluster, stateMachine.GetCluster("test2"))
	assert.Nil(t, stateMachine.GetCluster("test1"))

	// OnDelete
	cluster.EXPECT().Close()
	stateMachine.OnDelete("/test/data/test2")
	assert.Equal(t, 0, len(stateMachine.GetAllCluster()))

	// OnCreate
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(state.NewMockRepository(ctrl), nil)
	clusterFactory.EXPECT().newCluster(gomock.Any()).Return(cluster, nil)
	stateMachine.OnCreate("/test/data/test1", encoding.JSONMarshal(&models.StorageState{Name: "test1"}))

	stateMachine.Cleanup()

	cluster.EXPECT().Close()
	discovery1.EXPECT().Close()
	_ = stateMachine.Close()
}
