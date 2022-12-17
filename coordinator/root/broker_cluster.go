// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package root

import (
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator/discovery"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/logger"
	statepkg "github.com/lindb/lindb/pkg/state"
)

//go:generate mockgen -source=./broker_cluster.go -destination=./broker_cluster_mock.go -package=root

// BrokerCluster represents broker cluster controller,
// 1) discovery active node list in broker cluster
type BrokerCluster interface {
	// Start starts the state machine for broker state change.
	Start() error
	// GetState returns the current state of broker cluster.
	GetState() *models.BrokerState
	// Close closes broker cluster controller
	Close()
}

// brokerCluster implements BrokerCluster controller, root will maintain multi broker clustcr.
type brokerCluster struct {
	brokerRepo statepkg.Repository
	cfg        *config.BrokerCluster
	state      *models.BrokerState
	stateMgr   StateManager

	sm discovery.StateMachine

	logger *logger.Logger
}

// newBrokerCluster creates broker cluster controller, init active node list if exist node, must return a broker cluster instance.
func newBrokerCluster(
	cfg *config.BrokerCluster,
	stateMgr StateManager,
	repoFactory statepkg.RepositoryFactory,
) (BrokerCluster, error) {
	var brokerRepo statepkg.Repository
	brokerRepo, err := repoFactory.CreateBrokerRepo(cfg.Config)
	if err != nil {
		return nil, err
	}
	cluster := &brokerCluster{
		brokerRepo: brokerRepo,
		stateMgr:   stateMgr,
		state:      models.NewBrokerState(cfg.Config.Namespace),
		logger:     logger.GetLogger("Root", "Broker"),
	}
	cluster.logger.Info("init broker cluster success", logger.String("broker", cfg.Config.Namespace))
	return cluster, nil
}

// Start starts the state machine for broker state change.
func (c *brokerCluster) Start() error {
	sm, err := c.stateMgr.GetStateMachineFactory().
		createBrokerNodeStateMachine(c.cfg.Config.Namespace, discovery.NewFactory(c.brokerRepo))
	if err != nil {
		return err
	}
	c.sm = sm

	c.logger.Info("start broker cluster successfully", logger.String("broker", c.cfg.Config.Namespace))
	return nil
}

// GetState returns the current state of broker cluster.
func (c *brokerCluster) GetState() *models.BrokerState {
	return c.state
}

// Close stops watch, and cleanups broker cluster's metadata.
func (c *brokerCluster) Close() {
	c.logger.Info("close broker cluster state machine", logger.String("broker", c.cfg.Config.Namespace))
	if c.sm != nil {
		if err := c.sm.Close(); err != nil {
			c.logger.Error("close broker node state machine of broker cluster",
				logger.String("broker", c.cfg.Config.Namespace), logger.Error(err), logger.Stack())
		}
	}
	if err := c.brokerRepo.Close(); err != nil {
		c.logger.Error("close state repo of broker cluster",
			logger.String("broker", c.cfg.Config.Namespace), logger.Error(err), logger.Stack())
	}
}
