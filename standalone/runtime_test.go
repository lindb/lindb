package standalone

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/pkg/fileutil"
	"github.com/lindb/lindb/pkg/server"
	"github.com/lindb/lindb/pkg/state"
)

var testPath = "./test_data"
var defaultStandaloneConfig = config.Standalone{
	StorageBase: *config.NewDefaultStorageBase(),
	BrokerBase:  *config.NewDefaultBrokerBase(),
	Logging:     *config.NewDefaultLogging(),
	ETCD:        *config.NewDefaultETCD(),
	Monitor:     *config.NewDefaultMonitor(),
}

func TestRuntime_Run(t *testing.T) {
	defer func() {
		_ = fileutil.RemoveDir(testPath)
	}()
	cfg := defaultStandaloneConfig
	cfg.StorageBase.TSDB.Dir = testPath
	standalone := NewStandaloneRuntime("test-version", cfg)
	err := standalone.Run()
	assert.NoError(t, err)
	assert.Equal(t, server.Running, standalone.State())
	err = standalone.Stop()
	assert.NoError(t, err)
	assert.Equal(t, server.Terminated, standalone.State())
	assert.Equal(t, "standalone", standalone.Name())
}

func TestRuntime_Run_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	cfg := defaultStandaloneConfig
	cfg.StorageBase.TSDB.Dir = testPath
	standalone := NewStandaloneRuntime("test-version", cfg)
	s := standalone.(*runtime)
	storage := server.NewMockService(ctrl)
	s.storage = storage
	storage.EXPECT().Run().Return(fmt.Errorf("err"))
	err := standalone.Run()
	assert.Error(t, err)

	standalone = NewStandaloneRuntime("test-version", cfg)
	// restart etcd err
	err = standalone.Run()
	assert.Error(t, err)
	storage.EXPECT().Stop().Return(fmt.Errorf("err")).AnyTimes()
	err = s.Stop()
	assert.NoError(t, err)
}

func TestRuntime_runServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()
	cfg := defaultStandaloneConfig
	standalone := NewStandaloneRuntime("test-version", cfg)
	s := standalone.(*runtime)
	storage := server.NewMockService(ctrl)
	s.storage = storage
	broker := server.NewMockService(ctrl)
	s.broker = broker
	storage.EXPECT().Run().Return(nil).AnyTimes()
	broker.EXPECT().Run().Return(fmt.Errorf("err"))
	err := s.runServer()
	assert.Error(t, err)
	storage.EXPECT().Stop().Return(fmt.Errorf("err"))
	broker.EXPECT().Stop().Return(fmt.Errorf("err"))
	err = s.Stop()
	assert.NoError(t, err)
}

func TestRuntime_cleanupState(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer func() {
		ctrl.Finish()
		_ = fileutil.RemoveDir(testPath)
	}()

	cfg := defaultStandaloneConfig
	cfg.StorageBase.TSDB.Dir = testPath
	standalone := NewStandaloneRuntime("test-version", cfg)
	s := standalone.(*runtime)
	repoFactory := state.NewMockRepositoryFactory(ctrl)
	s.repoFactory = repoFactory
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(nil, fmt.Errorf("err"))
	err := standalone.Run()
	assert.Error(t, err)
	err = s.Stop()
	assert.NoError(t, err)

	repo := state.NewMockRepository(ctrl)
	repoFactory.EXPECT().CreateRepo(gomock.Any()).Return(repo, nil)
	repo.EXPECT().Delete(gomock.Any(), gomock.Any()).Return(fmt.Errorf("err"))
	repo.EXPECT().Close().Return(fmt.Errorf("err"))
	err = s.cleanupState()
	assert.Error(t, err)
}
