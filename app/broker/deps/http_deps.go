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

package deps

import (
	"context"

	"github.com/lindb/lindb/config"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/pkg/state"
	brokerQuery "github.com/lindb/lindb/query/broker"
	"github.com/lindb/lindb/replica"
)

// HTTPDeps represents http server handler's dependency.
type HTTPDeps struct {
	Ctx       context.Context
	BrokerCfg *config.Broker
	Master    coordinator.Master

	Repo     state.Repository
	StateMgr broker.StateManager

	CM replica.ChannelManager

	QueryFactory brokerQuery.Factory
}

func (deps *HTTPDeps) WithTimeout() (context.Context, context.CancelFunc) {
	// choose the shorter duration
	timeout := deps.BrokerCfg.Coordinator.Timeout.Duration()
	if deps.BrokerCfg.BrokerBase.HTTP.ReadTimeout.Duration() < timeout {
		timeout = deps.BrokerCfg.BrokerBase.HTTP.ReadTimeout.Duration()
	}
	return context.WithTimeout(deps.Ctx, timeout)
}
