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

package meta

import (
	"context"

	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	"github.com/lindb/lindb/spi/types"
)

type MetadataManager interface {
	GetStateRepo() state.Repository
	GetCurrentNode() models.StatelessNode
	// GetMaster returns the current master info.
	GetMaster() *models.Master
	// GetBrokerNodes returns all alive broker nodes.
	GetBrokerNodes() (nodes []models.StatelessNode)
	// GetStorageNodes returns all alive storage nodes.
	GetStorageNodes() (nodes []models.StatefulNode)
	// GetDatabaseCfg returns the database config by name.
	GetDatabase(database string) (models.Database, bool)
	// GetDatabases returns current database config list.
	GetDatabases() []models.Database
	GetPartitions(database, ns, table string) (partitions map[models.InternalNode][]int, err error)
	GetTableMetadata(database, ns, table string) (*types.TableMetadata, error)

	CreateDatabase(ctx context.Context, database *models.Database) error
	DropDatabase(ctx context.Context, database string) error
}
