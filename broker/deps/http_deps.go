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
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/parallel"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/replication"
	"github.com/lindb/lindb/service"
)

// HTTPDeps represents http server handler's dependency.
type HTTPDeps struct {
	Master coordinator.Master

	Repo          state.Repository
	StateMachines *coordinator.BrokerStateMachines

	DatabaseSrv       service.DatabaseService
	ShardAssignSrv    service.ShardAssignService
	StorageClusterSrv service.StorageClusterService

	CM replication.ChannelManager

	ExecutorFct parallel.ExecutorFactory
	JobManager  parallel.JobManager
}
