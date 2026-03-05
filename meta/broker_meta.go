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
	"fmt"
	"strings"

	"github.com/apache/arrow-go/v18/arrow"
	"github.com/lindb/common/pkg/encoding"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/models"
	"github.com/lindb/lindb/pkg/state"
	protoMetaV1 "github.com/lindb/lindb/proto/gen/v1/meta"
	"github.com/lindb/lindb/spi/types"
)

type brokerMetadataManager struct {
	repo             state.Repository
	brokerStateMgr   broker.StateManager
	masterStateMgr   master.StateManager
	masterController coordinator.MasterController
}

func NewBrokerMetadataManager(
	repo state.Repository,
	brokerStateMgr broker.StateManager,
	masterController coordinator.MasterController,
) MetadataManager {
	return &brokerMetadataManager{
		repo:             repo,
		brokerStateMgr:   brokerStateMgr,
		masterStateMgr:   masterController.GetStateManager(),
		masterController: masterController,
	}
}

func (m *brokerMetadataManager) GetStateRepo() state.Repository {
	return m.repo
}

func (m *brokerMetadataManager) GetCurrentNode() models.StatelessNode {
	return m.brokerStateMgr.GetCurrentNode()
}

func (m *brokerMetadataManager) GetMaster() *models.Master {
	return m.masterController.GetMaster()
}

func (m *brokerMetadataManager) GetBrokerNodes() (nodes []models.StatelessNode) {
	return m.brokerStateMgr.GetLiveNodes()
}

func (m *brokerMetadataManager) GetStorageNodes() (nodes []models.StatefulNode) {
	liveNodes := m.brokerStateMgr.GetStorage().LiveNodes
	for _, node := range liveNodes {
		nodes = append(nodes, node)
	}
	return
}

func (m *brokerMetadataManager) GetDatabase(database string) (models.Database, bool) {
	return m.brokerStateMgr.GetDatabase(database)
}

func (m *brokerMetadataManager) GetStreaming(streaming string) (models.Streaming, bool) {
	return m.brokerStateMgr.GetStreaming(streaming)
}

func (m *brokerMetadataManager) GetDatabases() []models.Database {
	return m.brokerStateMgr.GetDatabases()
}

func (m *brokerMetadataManager) GetPartitions(database, ns, table string) (map[models.InternalNode][]int, error) {
	if database == constants.InformationSchema {
		var partitions map[models.InternalNode][]int
		currentNode := m.brokerStateMgr.GetCurrentNode()
		switch strings.ToLower(table) {
		case constants.TableEngines,
			constants.TableEnv,
			constants.TableSchemata,
			constants.TableMetrics,
			constants.TableMaster,
			constants.TableBrokers,
			constants.TableStorages,
			constants.TableReplications,
			constants.TableMemoryDatabases,
			constants.TableNamespaces,
			constants.TableTableNames,
			constants.TableMetadataTypes,
			constants.TableMetadatas,
			constants.TableColumns,
			constants.TableFunctions,
			constants.TableSnippets,
			constants.TableStreamings,
			constants.TableStreamingJobs:
			partitions = map[models.InternalNode][]int{
				{IP: currentNode.HostIP, Port: currentNode.GRPCPort}: {},
			}
		}
		return partitions, nil
	}

	// find tabel metadata from partitions
	return m.brokerStateMgr.GetPartitions(database)
}

func (m *brokerMetadataManager) GetTableMetadata(database, ns, table string) (*types.TableMetadata, error) {
	// find tabel metadata from partitions
	partitions, err := m.GetPartitions(database, ns, table)
	if err != nil {
		return nil, err
	}
	var fields []arrow.Field
	supportDynamicField := false
	if table != "logs" && table != "traces" {
		// FIXME: log table???
		for node := range partitions {
			tableSchema, err := m.getTableSchema(database, ns, table, node)
			if err != nil {
				return nil, err
			}
			// TODO: remove duplicate column
			fields = append(fields, tableSchema.Fields()...)
		}
	} else {
		supportDynamicField = true // for log/tarce
		// TODO: add hidden flag?
		fields = append(fields, arrow.Field{Name: constants.TimestampColumnName, Type: arrow.FixedWidthTypes.Timestamp_s})
	}
	return &types.TableMetadata{
		Schema:              arrow.NewSchema(fields, nil),
		Partitions:          partitions,
		SupportDynamicField: supportDynamicField,
	}, nil
}

func (m *brokerMetadataManager) CreateDatabase(ctx context.Context, database *models.Database) error {
	return m.repo.Put(ctx, constants.GetDatabaseConfigPath(database.Name), encoding.JSONMarshal(database))
}

func (m *brokerMetadataManager) DropDatabase(ctx context.Context, database string) error {
	return m.repo.Delete(ctx, constants.GetDatabaseConfigPath(database))
}

func (m *brokerMetadataManager) CreateStreaming(ctx context.Context, streaming *models.Streaming) error {
	return m.repo.Put(ctx, constants.GetStreamingConfigPath(streaming.Name), encoding.JSONMarshal(streaming))
}

func (m *brokerMetadataManager) DropStreaming(ctx context.Context, streaming string) error {
	if err := m.repo.Delete(ctx, constants.GetStreamingStatePath(streaming)); err != nil {
		return nil
	}
	if err := m.repo.Delete(ctx, constants.GetStreamingJobsPath(streaming)); err != nil {
		return nil
	}
	return m.repo.Delete(ctx, constants.GetStreamingConfigPath(streaming))
}

func (m *brokerMetadataManager) CreateJob(ctx context.Context, stream, job, sql string) error {
	jobPath := constants.GetStreamingJobPath(stream, job)
	return m.repo.Put(ctx, jobPath, []byte(sql))
}

func (m *brokerMetadataManager) DropJob(ctx context.Context, stream, job string) error {
	jobPath := constants.GetStreamingJobPath(stream, job)
	return m.repo.Delete(ctx, jobPath)
}

func (m *brokerMetadataManager) getTableSchema(
	database, ns, table string,
	node models.InternalNode,
) (*arrow.Schema, error) {
	conn, err := grpc.NewClient(node.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := protoMetaV1.NewMetaServiceClient(conn)
	resp, err := client.TableSchema(context.TODO(), &protoMetaV1.TableSchemaRequest{
		Database:  database,
		Namespace: ns,
		Table:     table,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)
	panic("FIXME: return table schema")
}
