package storage

import (
	"context"

	"github.com/eleme/lindb/config"
	"github.com/eleme/lindb/coordinator/discovery"
	"github.com/eleme/lindb/pkg/logger"

	"go.uber.org/zap"
)

type State uint8

const (
	Created State = iota
	Running State = 1
)

type Storage struct {
	state    State
	ctx      context.Context
	register *discovery.Register
	config   *config.StorageConfig
	cancel   context.CancelFunc
	log      *zap.Logger
}

// New returns a storage
func New(ctx context.Context, storageConfig *config.StorageConfig) *Storage {
	ctx, cancel := context.WithCancel(ctx)
	return &Storage{
		ctx:    ctx,
		config: storageConfig,
		cancel: cancel,
		state:  Created,
		log:    logger.GetLogger(),
	}
}

// Start starts the storage and  registers the storage to the repository
func (s *Storage) Start() error {
	if err := s.initRepository(); err != nil {
		return err
	}
	// ip, err := util.GetHostIP()
	// if err != nil {
	// 	return err
	// }
	// nodeName := fmt.Sprintf("%s:%v", ip, s.config.StoragePort)
	// key := fmt.Sprintf("%s/storage/%s/", s.config.HeartBeatPrefix, nodeName)
	// node := models.Node{IP: ip, Port: s.config.StoragePort}
	// s.register = discovery.NewRegister(key, node, s.config.HeartBeatTTL)
	// if err := s.register.Register(s.ctx); err != nil {
	// 	return err
	// }
	// //TODO gets cluster config from etcd and registers to the cluster
	// s.state = Running
	return nil
}

// initRepository inits the repository for the storage
func (s *Storage) initRepository() error {
	// newRepositoryConfig, err := state.NewRepositoryConfig(s.config.RepositoryConfig)
	// if err != nil {
	// 	s.log.Error("convent repository config failed", zap.Error(err))
	// 	return err
	// }
	// if err := state.New(s.config.RepositoryType, newRepositoryConfig); err != nil {
	// 	s.log.Error("init repository failed", zap.Error(err))
	// 	return err
	// }
	return nil
}

// Close closes storage server
func (s *Storage) Close() error {
	if s.state != Running {
		s.log.Warn("the storage is not running,needn't stop")
		return nil
	}
	if err := s.register.UnRegister(s.ctx); err != nil {
		s.log.Error("unregister storage error", zap.Error(err))
	}
	s.cancel()
	return nil
}
