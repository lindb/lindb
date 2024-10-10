package meta

import (
	"strings"

	"github.com/lindb/lindb/constants"
	"github.com/lindb/lindb/coordinator"
	"github.com/lindb/lindb/coordinator/broker"
	"github.com/lindb/lindb/coordinator/master"
	"github.com/lindb/lindb/models"
)

type brokerMetadataManager struct {
	brokerStateMgr   broker.StateManager
	masterStateMgr   master.StateManager
	masterController coordinator.MasterController
}

func NewBrokerMetadataManager(
	brokerStateMgr broker.StateManager,
	masterController coordinator.MasterController,
) MetadataManager {
	return &brokerMetadataManager{
		brokerStateMgr:   brokerStateMgr,
		masterStateMgr:   masterController.GetStateManager(),
		masterController: masterController,
	}
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
		case constants.TableEngines, constants.TableSchemata, constants.TableMetrics, constants.TableMaster, constants.TableBroker, constants.TableStorage:
			partitions = map[models.InternalNode][]int{
				{IP: currentNode.HostIP, Port: currentNode.GRPCPort}: {},
			}
		}
		return partitions, nil
	}

	// find tabel metadata from partitions
	return m.brokerStateMgr.GetPartitions(database)
}
