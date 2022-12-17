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

package command

import (
	"bytes"

	"context"
	"errors"
	"fmt"
	"time"

	depspkg "github.com/lindb/lindb/app/root/deps"
	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/encoding"
	"github.com/lindb/lindb/pkg/logger"
	"github.com/lindb/lindb/pkg/ltoml"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/pkg/validate"
	stmtpkg "github.com/lindb/lindb/sql/stmt"
)

var log = logger.GetLogger("Exec", "Command")

// brokerCommandFn represents broker command function define.
type brokerCommandFn = func(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Broker) (interface{}, error)

// brokerCommands registers all broker related commands.
var brokerCommands = map[stmtpkg.BrokerOpType]brokerCommandFn{
	stmtpkg.BrokerOpShow:   listBrokers,
	stmtpkg.BrokerOpCreate: createBroker,
}

// BrokerCommand executes lin query language for broker related.
func BrokerCommand(ctx context.Context, deps *depspkg.HTTPDeps, _ *models.ExecuteParam, stmt stmtpkg.Statement) (interface{}, error) {
	brokerStmt := stmt.(*stmtpkg.Broker)
	if commandFn, ok := brokerCommands[brokerStmt.Type]; ok {
		return commandFn(ctx, deps, brokerStmt)
	}
	return nil, nil
}

// List lists all broker clusters
func listBrokers(ctx context.Context, deps *depspkg.HTTPDeps, _ *stmtpkg.Broker) (interface{}, error) {
	data, err := deps.Repo.List(ctx, constants.StorageConfigPath)
	if err != nil {
		return nil, err
	}
	stateMgr := deps.StateMgr
	var brokers models.Brokers
	for _, val := range data {
		broker := models.Broker{}
		err = encoding.JSONUnmarshal(val.Value, &broker)
		if err != nil {
			log.Warn("unmarshal data error",
				logger.String("data", string(val.Value)))
		} else {
			if _, ok := stateMgr.GetBrokerState(broker.Config.Namespace); ok {
				broker.Status = models.ClusterStatusReady
			} else {
				broker.Status = models.ClusterStatusInitialize
				// TODO: check broker un-health
			}
			brokers = append(brokers, broker)
		}
	}

	if err != nil {
		return nil, err
	}
	return brokers, nil
}

// createBroker creates config of broker cluster.
func createBroker(ctx context.Context, deps *depspkg.HTTPDeps, stmt *stmtpkg.Broker) (interface{}, error) {
	data := []byte(stmt.Value)
	broker := &config.BrokerCluster{}
	err := encoding.JSONUnmarshal(data, broker)
	if err != nil {
		return nil, err
	}
	err = validate.Validator.Struct(broker)
	if err != nil {
		return nil, err
	}
	// copy config for testing
	cfg := &config.RepoState{}
	_ = encoding.JSONUnmarshal(encoding.JSONMarshal(broker.Config), cfg)
	cfg.Timeout = ltoml.Duration(time.Second)
	cfg.DialTimeout = ltoml.Duration(time.Second)
	// check broker repo config if valid
	repo, err := deps.RepoFactory.CreateBrokerRepo(cfg)
	if err != nil {
		return nil, err
	}
	err = repo.Close()
	if err != nil {
		return nil, err
	}
	// re-marshal broker config, keep same structure with repo.
	data = encoding.JSONMarshal(broker)
	log.Info("Creating broker cluster", logger.String("config", stmt.Value))
	ok, err := deps.Repo.PutWithTX(ctx, constants.GetBrokerClusterConfigPath(broker.Config.Namespace), data, func(oldVal []byte) error {
		if bytes.Equal(data, oldVal) {
			log.Info("broker cluster exist", logger.String("config", string(oldVal)))
			return state.ErrNotExist
		}
		return nil
	})
	if errors.Is(state.ErrNotExist, err) {
		rs := "Broker is exist"
		return &rs, nil
	}
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("create broker failure")
	}
	rs := "Create broker ok"
	return &rs, nil
}
