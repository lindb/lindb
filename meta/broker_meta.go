package meta

import (
	"context"
	"strings"

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

func (m *brokerMetadataManager) GetDatabases() []models.Database {
	return m.brokerStateMgr.GetDatabases()
}

func (m *brokerMetadataManager) GetPartitions(database, ns, table string) (map[models.InternalNode][]int, error) {
	if database == constants.InformationSchema {
		var partitions map[models.InternalNode][]int
		currentNode := m.brokerStateMgr.GetCurrentNode()
		switch strings.ToLower(table) {
		case constants.TableEngines,
			constants.TableSchemata,
			constants.TableMetrics,
			constants.TableMaster,
			constants.TableBroker,
			constants.TableStorage,
			constants.TableReplications,
			constants.TableMemoryDatabases,
			constants.TableNamespaces,
			constants.TableTableNames,
			constants.TableMetadataTypes,
			constants.TableMetadatas,
			constants.TableColumns:
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
	schema := types.NewTableSchema()
	for node := range partitions {
		conn, err := grpc.Dial(node.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
		tableSchema := &types.TableSchema{}
		if err = encoding.JSONUnmarshal(resp.Payload, tableSchema); err != nil {
			return nil, err
		}
		// TODO: remove duplicate column
		schema.AddColumns(tableSchema.Columns)
	}
	return &types.TableMetadata{
		Schema:     schema,
		Partitions: partitions,
	}, nil
}
